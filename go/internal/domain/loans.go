package domain

import (
	"fmt"
	"strings"

	"mini-world-go/internal/world"
)

type LoanRepaymentResult struct {
	InterestPaid  int
	PrincipalPaid int
	AccountID     string
}

func GrantLoan(w *world.World, bankID string, humanID string, currency string, amount int, collateralAssetID string) (string, string, error) {
	if amount <= 0 {
		return "", "", fmt.Errorf("amount must be greater than 0")
	}

	bank, bankExists := w.Banks[bankID]
	if !bankExists {
		return "", "", fmt.Errorf("bank does not exist: %s", bankID)
	}

	human, humanExists := w.Humans[humanID]
	if !humanExists {
		return "", "", fmt.Errorf("human does not exist: %s", humanID)
	}

	currency = strings.ToUpper(currency)
	if _, exists := w.Currencies[currency]; !exists {
		return "", "", fmt.Errorf("currency does not exist: %s", currency)
	}

	accountID, account, err := activeAccount(w, humanID, bankID, currency)
	if err != nil {
		return "", "", err
	}

	if collateralAssetID != "" {
		collateral, ok := w.Assets[collateralAssetID]
		if !ok {
			return "", "", fmt.Errorf("collateral asset does not exist: %s", collateralAssetID)
		}

		if collateral.OwnerType != world.AssetOwnerHuman || collateral.OwnerID != humanID {
			return "", "", fmt.Errorf("collateral asset is not owned by human: %s", humanID)
		}

		if collateral.Currency != currency {
			return "", "", fmt.Errorf("collateral asset currency is %s, not %s", collateral.Currency, currency)
		}

		if collateral.CollateralForReserveLoanID != "" {
			return "", "", fmt.Errorf("collateral asset is already pledged to reserve loan: %s", collateral.CollateralForReserveLoanID)
		}

		if collateral.CollateralForCustomerLoanID != "" {
			return "", "", fmt.Errorf("collateral asset is already pledged to customer loan: %s", collateral.CollateralForCustomerLoanID)
		}

		if collateral.EstimatedValue < amount {
			return "", "", fmt.Errorf(
				"collateral value is too low. Estimated value: %d %s",
				collateral.EstimatedValue,
				currency,
			)
		}
	}

	loanID := nextCustomerLoanID(w)
	w.CustomerLoans[loanID] = world.NewCustomerLoan(loanID, bankID, humanID, currency, amount)
	w.CustomerLoans[loanID].CollateralAssetID = collateralAssetID
	bank.CustomerLoans = append(bank.CustomerLoans, loanID)
	human.Loans[loanID] = amount
	account.BookedBalance += amount

	if collateralAssetID != "" {
		collateral := w.Assets[collateralAssetID]
		collateral.PledgedToBankID = bankID
		collateral.CollateralForCustomerLoanID = loanID
	}

	return loanID, accountID, nil
}

func AccrueInterest(w *world.World, loanID string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	loan, ok := w.CustomerLoans[loanID]
	if !ok {
		return fmt.Errorf("loan does not exist: %s", loanID)
	}

	if loan.Status != world.CustomerLoanActive {
		return fmt.Errorf("loan is not active: %s", loanID)
	}

	if _, exists := w.Banks[loan.BankID]; !exists {
		return fmt.Errorf("bank does not exist: %s", loan.BankID)
	}

	human, humanExists := w.Humans[loan.BorrowerHumanID]
	if !humanExists {
		return fmt.Errorf("human does not exist: %s", loan.BorrowerHumanID)
	}

	loan.OutstandingInterest += amount
	loan.TotalInterestAccrued += amount
	human.Loans[loanID] = loan.OutstandingPrincipal + loan.OutstandingInterest

	return nil
}

func RepayLoan(w *world.World, humanID string, bankID string, loanID string, amount int) (LoanRepaymentResult, error) {
	if amount <= 0 {
		return LoanRepaymentResult{}, fmt.Errorf("amount must be greater than 0")
	}

	human, humanExists := w.Humans[humanID]
	if !humanExists {
		return LoanRepaymentResult{}, fmt.Errorf("human does not exist: %s", humanID)
	}

	bank, bankExists := w.Banks[bankID]
	if !bankExists {
		return LoanRepaymentResult{}, fmt.Errorf("bank does not exist: %s", bankID)
	}

	loan, ok := w.CustomerLoans[loanID]
	if !ok {
		return LoanRepaymentResult{}, fmt.Errorf("loan does not exist: %s", loanID)
	}

	if loan.BankID != bankID {
		return LoanRepaymentResult{}, fmt.Errorf("loan %s does not belong to %s", loanID, bankID)
	}

	if loan.BorrowerHumanID != humanID {
		return LoanRepaymentResult{}, fmt.Errorf("loan %s does not belong to %s", loanID, humanID)
	}

	if loan.Status != world.CustomerLoanActive {
		return LoanRepaymentResult{}, fmt.Errorf("loan is not active: %s", loanID)
	}

	totalDue := loan.OutstandingPrincipal + loan.OutstandingInterest
	if totalDue < amount {
		return LoanRepaymentResult{}, fmt.Errorf("repayment amount is greater than total due. Total due: %d %s", totalDue, loan.Currency)
	}

	accountID, account, err := activeAccount(w, humanID, bankID, loan.Currency)
	if err != nil {
		return LoanRepaymentResult{}, err
	}

	if account.BookedBalance < amount {
		return LoanRepaymentResult{}, fmt.Errorf(
			"not enough money in %s's account. Available: %d %s",
			humanID,
			account.BookedBalance,
			loan.Currency,
		)
	}

	account.BookedBalance -= amount
	remainingPayment := amount

	interestPaid := min(remainingPayment, loan.OutstandingInterest)
	loan.OutstandingInterest -= interestPaid
	remainingPayment -= interestPaid

	principalPaid := min(remainingPayment, loan.OutstandingPrincipal)
	loan.OutstandingPrincipal -= principalPaid

	bank.InterestIncome[loan.Currency] += interestPaid
	bank.Equity[loan.Currency] += interestPaid
	human.Loans[loanID] = loan.OutstandingPrincipal + loan.OutstandingInterest

	if human.Loans[loanID] == 0 {
		loan.Status = world.CustomerLoanRepaid
		if loan.CollateralAssetID != "" {
			if collateral, ok := w.Assets[loan.CollateralAssetID]; ok {
				collateral.PledgedToBankID = ""
				collateral.CollateralForCustomerLoanID = ""
			}
		}
	}

	return LoanRepaymentResult{
		InterestPaid:  interestPaid,
		PrincipalPaid: principalPaid,
		AccountID:     accountID,
	}, nil
}

func DefaultLoan(w *world.World, loanID string) error {
	loan, ok := w.CustomerLoans[loanID]
	if !ok {
		return fmt.Errorf("loan does not exist: %s", loanID)
	}

	if loan.Status == world.CustomerLoanRepaid {
		return fmt.Errorf("loan is already repaid: %s", loanID)
	}

	if loan.Status == world.CustomerLoanDefaulted {
		return fmt.Errorf("loan is already defaulted: %s", loanID)
	}

	if loan.Status != world.CustomerLoanActive {
		return fmt.Errorf("loan is not active: %s", loanID)
	}

	loan.Status = world.CustomerLoanDefaulted

	return nil
}

func nextCustomerLoanID(w *world.World) string {
	for index := 1; ; index++ {
		loanID := fmt.Sprintf("loan_%06d", index)
		if _, exists := w.CustomerLoans[loanID]; !exists {
			return loanID
		}
	}
}

func min(left int, right int) int {
	if left < right {
		return left
	}

	return right
}
