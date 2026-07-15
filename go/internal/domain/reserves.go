package domain

import (
	"fmt"
	"strings"

	"mini-world-go/internal/world"
)

func LendReserves(w *world.World, centralBankID string, bankID string, currency string, amount int, collateralAssetID string) (string, error) {
	if amount <= 0 {
		return "", fmt.Errorf("amount must be greater than 0")
	}

	centralBank, ok := w.CentralBanks[centralBankID]
	if !ok {
		return "", fmt.Errorf("central bank does not exist: %s", centralBankID)
	}

	bank, ok := w.Banks[bankID]
	if !ok {
		return "", fmt.Errorf("bank does not exist: %s", bankID)
	}

	currency = strings.ToUpper(currency)
	if centralBank.Currency != currency {
		return "", fmt.Errorf("%s issues %s, not %s", centralBankID, centralBank.Currency, currency)
	}

	if _, err := reserveBalance(centralBank, bank, centralBankID, bankID); err != nil {
		return "", err
	}

	if collateralAssetID != "" {
		collateral, ok := w.Assets[collateralAssetID]
		if !ok {
			return "", fmt.Errorf("collateral asset does not exist: %s", collateralAssetID)
		}

		if collateral.OwnerType != world.AssetOwnerBank || collateral.OwnerID != bankID {
			return "", fmt.Errorf("collateral asset is not owned by bank: %s", bankID)
		}

		if collateral.Currency != currency {
			return "", fmt.Errorf("collateral asset currency is %s, not %s", collateral.Currency, currency)
		}

		if collateral.CollateralForReserveLoanID != "" {
			return "", fmt.Errorf("collateral asset is already pledged to reserve loan: %s", collateral.CollateralForReserveLoanID)
		}

		if collateral.EstimatedValue < amount {
			return "", fmt.Errorf(
				"collateral value is too low. Estimated value: %d %s",
				collateral.EstimatedValue,
				currency,
			)
		}
	}

	reserveLoanID := nextReserveLoanID(w, centralBankID, bankID)
	w.ReserveLoans[reserveLoanID] = world.NewReserveLoan(reserveLoanID, centralBankID, bankID, currency, amount, collateralAssetID)

	centralBank.ReserveAccounts[bankID] += amount
	bank.ReserveBalances[centralBankID] += amount
	centralBank.LoansToBanks[bankID] += amount
	bank.LoansFromCentralBanks[centralBankID] += amount
	if collateralAssetID != "" {
		collateral := w.Assets[collateralAssetID]
		collateral.PledgedToCentralBankID = centralBankID
		collateral.CollateralForReserveLoanID = reserveLoanID
	}

	return reserveLoanID, nil
}

func RepayReserveLoan(w *world.World, bankID string, reserveLoanID string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	bank, ok := w.Banks[bankID]
	if !ok {
		return fmt.Errorf("bank does not exist: %s", bankID)
	}

	reserveLoan, ok := w.ReserveLoans[reserveLoanID]
	if !ok {
		return fmt.Errorf("reserve loan does not exist: %s", reserveLoanID)
	}

	if reserveLoan.BankID != bankID {
		return fmt.Errorf("reserve loan %s belongs to %s, not %s", reserveLoanID, reserveLoan.BankID, bankID)
	}

	if reserveLoan.Status != world.ReserveLoanOpen {
		return fmt.Errorf("reserve loan is not open: %s", reserveLoanID)
	}

	if reserveLoan.Outstanding < amount {
		return fmt.Errorf(
			"%s does not owe that much. Outstanding loan: %d %s",
			bankID,
			reserveLoan.Outstanding,
			reserveLoan.Currency,
		)
	}

	centralBank, ok := w.CentralBanks[reserveLoan.CentralBankID]
	if !ok {
		return fmt.Errorf("central bank does not exist: %s", reserveLoan.CentralBankID)
	}

	reserves, err := reserveBalance(centralBank, bank, reserveLoan.CentralBankID, bankID)
	if err != nil {
		return err
	}

	if reserves < amount {
		return fmt.Errorf(
			"not enough reserves for %s at %s. Available: %d %s",
			bankID,
			reserveLoan.CentralBankID,
			reserves,
			reserveLoan.Currency,
		)
	}

	if err := ensureLoanMirror(centralBank, bank, reserveLoan.CentralBankID, bankID); err != nil {
		return err
	}

	centralBank.ReserveAccounts[bankID] -= amount
	bank.ReserveBalances[reserveLoan.CentralBankID] -= amount
	centralBank.LoansToBanks[bankID] -= amount
	bank.LoansFromCentralBanks[reserveLoan.CentralBankID] -= amount
	reserveLoan.Outstanding -= amount

	if reserveLoan.Outstanding == 0 {
		reserveLoan.Status = world.ReserveLoanRepaid
		if collateral, ok := w.Assets[reserveLoan.CollateralAssetID]; ok {
			collateral.PledgedToCentralBankID = ""
			collateral.CollateralForReserveLoanID = ""
		}
	}

	return nil
}

func ReserveTransfer(w *world.World, centralBankID string, fromBankID string, toBankID string, currency string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if fromBankID == toBankID {
		return fmt.Errorf("from bank and to bank must be different")
	}

	centralBank, ok := w.CentralBanks[centralBankID]
	if !ok {
		return fmt.Errorf("central bank does not exist: %s", centralBankID)
	}

	fromBank, ok := w.Banks[fromBankID]
	if !ok {
		return fmt.Errorf("from bank does not exist: %s", fromBankID)
	}

	toBank, ok := w.Banks[toBankID]
	if !ok {
		return fmt.Errorf("to bank does not exist: %s", toBankID)
	}

	currency = strings.ToUpper(currency)
	if centralBank.Currency != currency {
		return fmt.Errorf("%s issues %s, not %s", centralBankID, centralBank.Currency, currency)
	}

	fromReserves, err := reserveBalance(centralBank, fromBank, centralBankID, fromBankID)
	if err != nil {
		return err
	}

	if _, err := reserveBalance(centralBank, toBank, centralBankID, toBankID); err != nil {
		return err
	}

	if fromReserves < amount {
		return fmt.Errorf(
			"not enough reserves for %s at %s. Available: %d %s",
			fromBankID,
			centralBankID,
			fromReserves,
			currency,
		)
	}

	centralBank.ReserveAccounts[fromBankID] -= amount
	fromBank.ReserveBalances[centralBankID] -= amount
	centralBank.ReserveAccounts[toBankID] += amount
	toBank.ReserveBalances[centralBankID] += amount

	return nil
}

func nextReserveLoanID(w *world.World, centralBankID string, bankID string) string {
	for index := 1; ; index++ {
		reserveLoanID := fmt.Sprintf("reserve_loan_%s_%s_%d", centralBankID, bankID, index)
		if _, exists := w.ReserveLoans[reserveLoanID]; !exists {
			return reserveLoanID
		}
	}
}

func ensureLoanMirror(centralBank *world.CentralBank, bank *world.Bank, centralBankID string, bankID string) error {
	centralBankLoan := centralBank.LoansToBanks[bankID]
	bankLoanMirror := bank.LoansFromCentralBanks[centralBankID]

	if centralBankLoan != bankLoanMirror {
		return fmt.Errorf(
			"loan mirror mismatch. %s.loans_to_banks[%s] = %d, %s.loans_from_central_banks[%s] = %d",
			centralBankID,
			bankID,
			centralBankLoan,
			bankID,
			centralBankID,
			bankLoanMirror,
		)
	}

	return nil
}
