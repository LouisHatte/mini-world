package domain

import (
	"fmt"

	"mini-world-go/internal/world"
)

func CheckWorld(w *world.World) []string {
	var errors []string
	errors = append(errors, checkReserveMirrors(w)...)
	errors = append(errors, checkLoanMirrors(w)...)
	errors = append(errors, checkCustomerLoanMirrors(w)...)
	errors = append(errors, checkCorrespondentAccountMirrors(w)...)
	return errors
}

func checkReserveMirrors(w *world.World) []string {
	var errors []string

	for centralBankID, centralBank := range w.CentralBanks {
		for bankID, reserveAmount := range centralBank.ReserveAccounts {
			bank, ok := w.Banks[bankID]
			if !ok {
				errors = append(errors, fmt.Sprintf("%s.reserve_accounts contains unknown bank: %s", centralBankID, bankID))
				continue
			}

			mirroredAmount, ok := bank.ReserveBalances[centralBankID]
			if !ok || mirroredAmount != reserveAmount {
				errors = append(errors, fmt.Sprintf(
					"Reserve mismatch: %s.reserve_accounts[%s] = %d, but %s.reserve_balances[%s] = %d",
					centralBankID,
					bankID,
					reserveAmount,
					bankID,
					centralBankID,
					mirroredAmount,
				))
			}
		}
	}

	for bankID, bank := range w.Banks {
		for centralBankID, mirroredAmount := range bank.ReserveBalances {
			centralBank, ok := w.CentralBanks[centralBankID]
			if !ok {
				errors = append(errors, fmt.Sprintf("%s.reserve_balances contains unknown central bank: %s", bankID, centralBankID))
				continue
			}

			reserveAmount, ok := centralBank.ReserveAccounts[bankID]
			if !ok || reserveAmount != mirroredAmount {
				errors = append(errors, fmt.Sprintf(
					"Reserve mismatch: %s.reserve_balances[%s] = %d, but %s.reserve_accounts[%s] = %d",
					bankID,
					centralBankID,
					mirroredAmount,
					centralBankID,
					bankID,
					reserveAmount,
				))
			}
		}
	}

	return errors
}

func checkLoanMirrors(w *world.World) []string {
	var errors []string

	for centralBankID, centralBank := range w.CentralBanks {
		for bankID, loanAmount := range centralBank.LoansToBanks {
			bank, ok := w.Banks[bankID]
			if !ok {
				errors = append(errors, fmt.Sprintf("%s.loans_to_banks contains unknown bank: %s", centralBankID, bankID))
				continue
			}

			mirroredAmount, ok := bank.LoansFromCentralBanks[centralBankID]
			if !ok || mirroredAmount != loanAmount {
				errors = append(errors, fmt.Sprintf(
					"Loan mismatch: %s.loans_to_banks[%s] = %d, but %s.loans_from_central_banks[%s] = %d",
					centralBankID,
					bankID,
					loanAmount,
					bankID,
					centralBankID,
					mirroredAmount,
				))
			}
		}
	}

	for bankID, bank := range w.Banks {
		for centralBankID, mirroredAmount := range bank.LoansFromCentralBanks {
			centralBank, ok := w.CentralBanks[centralBankID]
			if !ok {
				errors = append(errors, fmt.Sprintf("%s.loans_from_central_banks contains unknown central bank: %s", bankID, centralBankID))
				continue
			}

			loanAmount, ok := centralBank.LoansToBanks[bankID]
			if !ok || loanAmount != mirroredAmount {
				errors = append(errors, fmt.Sprintf(
					"Loan mismatch: %s.loans_from_central_banks[%s] = %d, but %s.loans_to_banks[%s] = %d",
					bankID,
					centralBankID,
					mirroredAmount,
					centralBankID,
					bankID,
					loanAmount,
				))
			}
		}
	}

	return errors
}

func checkCustomerLoanMirrors(w *world.World) []string {
	var errors []string

	for loanID, loan := range w.CustomerLoans {
		bank, bankExists := w.Banks[loan.BankID]
		if !bankExists {
			errors = append(errors, fmt.Sprintf("%s references unknown bank: %s", loanID, loan.BankID))
		} else if !containsString(bank.CustomerLoans, loanID) {
			errors = append(errors, fmt.Sprintf("%s.customer_loans does not contain %s", loan.BankID, loanID))
		}

		human, humanExists := w.Humans[loan.BorrowerHumanID]
		if !humanExists {
			errors = append(errors, fmt.Sprintf("%s references unknown human: %s", loanID, loan.BorrowerHumanID))
			continue
		}

		expectedAmount := loan.OutstandingPrincipal + loan.OutstandingInterest
		mirroredAmount, ok := human.Loans[loanID]
		if !ok || mirroredAmount != expectedAmount {
			errors = append(errors, fmt.Sprintf(
				"Customer loan mismatch: %s expected due = %d, but %s.loans[%s] = %d",
				loanID,
				expectedAmount,
				loan.BorrowerHumanID,
				loanID,
				mirroredAmount,
			))
		}
	}

	for humanID, human := range w.Humans {
		for loanID, mirroredAmount := range human.Loans {
			loan, ok := w.CustomerLoans[loanID]
			if !ok {
				errors = append(errors, fmt.Sprintf("%s.loans contains unknown loan: %s", humanID, loanID))
				continue
			}

			expectedAmount := loan.OutstandingPrincipal + loan.OutstandingInterest
			if loan.BorrowerHumanID != humanID || mirroredAmount != expectedAmount {
				errors = append(errors, fmt.Sprintf(
					"Customer loan mismatch: %s.loans[%s] = %d, but %s borrower = %s and due = %d",
					humanID,
					loanID,
					mirroredAmount,
					loanID,
					loan.BorrowerHumanID,
					expectedAmount,
				))
			}
		}
	}

	return errors
}

func checkCorrespondentAccountMirrors(w *world.World) []string {
	var errors []string

	for correspondentAccountID, correspondentAccount := range w.CorrespondentAccounts {
		ownerBank, ownerBankExists := w.Banks[correspondentAccount.OwnerBankID]
		if !ownerBankExists {
			errors = append(errors, fmt.Sprintf("%s references unknown owner bank: %s", correspondentAccountID, correspondentAccount.OwnerBankID))
		} else if !containsString(ownerBank.NostroAccounts, correspondentAccountID) {
			errors = append(errors, fmt.Sprintf("%s.nostro_accounts does not contain %s", correspondentAccount.OwnerBankID, correspondentAccountID))
		}

		correspondentBank, correspondentBankExists := w.Banks[correspondentAccount.CorrespondentBankID]
		if !correspondentBankExists {
			errors = append(errors, fmt.Sprintf("%s references unknown correspondent bank: %s", correspondentAccountID, correspondentAccount.CorrespondentBankID))
		} else if !containsString(correspondentBank.VostroAccounts, correspondentAccountID) {
			errors = append(errors, fmt.Sprintf("%s.vostro_accounts does not contain %s", correspondentAccount.CorrespondentBankID, correspondentAccountID))
		}
	}

	for bankID, bank := range w.Banks {
		for _, correspondentAccountID := range bank.NostroAccounts {
			correspondentAccount, ok := w.CorrespondentAccounts[correspondentAccountID]
			if !ok {
				errors = append(errors, fmt.Sprintf("%s.nostro_accounts contains unknown correspondent account: %s", bankID, correspondentAccountID))
				continue
			}

			if correspondentAccount.OwnerBankID != bankID {
				errors = append(errors, fmt.Sprintf("%s.nostro_accounts contains %s, but owner bank is %s", bankID, correspondentAccountID, correspondentAccount.OwnerBankID))
			}
		}

		for _, correspondentAccountID := range bank.VostroAccounts {
			correspondentAccount, ok := w.CorrespondentAccounts[correspondentAccountID]
			if !ok {
				errors = append(errors, fmt.Sprintf("%s.vostro_accounts contains unknown correspondent account: %s", bankID, correspondentAccountID))
				continue
			}

			if correspondentAccount.CorrespondentBankID != bankID {
				errors = append(errors, fmt.Sprintf("%s.vostro_accounts contains %s, but correspondent bank is %s", bankID, correspondentAccountID, correspondentAccount.CorrespondentBankID))
			}
		}
	}

	return errors
}

func containsString(values []string, searchedValue string) bool {
	for _, value := range values {
		if value == searchedValue {
			return true
		}
	}

	return false
}
