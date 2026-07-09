package domain

import (
	"fmt"

	"mini-world-go/internal/world"
)

func CheckWorld(w *world.World) []string {
	var errors []string
	errors = append(errors, checkReserveMirrors(w)...)
	errors = append(errors, checkLoanMirrors(w)...)
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
