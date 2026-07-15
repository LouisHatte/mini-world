package loans

import (
	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newRepayLoanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "repay-loan human_id bank_id loan_id amount",
		Short: "Repay a commercial bank loan using the borrower's deposit balance.",
		Args:  cobra.ExactArgs(4),
		RunE:  runRepayLoan,
	}
}

func runRepayLoan(cmd *cobra.Command, args []string) error {
	humanID := args[0]
	bankID := args[1]
	loanID := args[2]
	amount, err := parseAmount(args[3])
	if err != nil {
		return err
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	currency := ""
	if loan, ok := w.CustomerLoans[loanID]; ok {
		currency = loan.Currency
	}

	result, err := domain.RepayLoan(w, humanID, bankID, loanID, amount)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "repay-loan", args); err != nil {
		return err
	}

	commandlog.Action("%s repaid %d %s on %s", humanID, amount, currency, loanID)
	commandlog.State("Interest paid: %d %s", result.InterestPaid, currency)
	commandlog.State("Principal paid: %d %s", result.PrincipalPaid, currency)
	commandlog.State("%s booked_balance: %d %s", result.AccountID, w.Accounts[result.AccountID].BookedBalance, currency)
	logLoanState(w, loanID)
	if collateralAssetID := w.CustomerLoans[loanID].CollateralAssetID; collateralAssetID != "" {
		commandlog.State("%s pledged_to_bank_id: %s", collateralAssetID, w.Assets[collateralAssetID].PledgedToBankID)
		commandlog.State("%s collateral_for_customer_loan_id: %s", collateralAssetID, w.Assets[collateralAssetID].CollateralForCustomerLoanID)
	}
	commandlog.State("%s interest_income[%s]: %d %s", bankID, currency, w.Banks[bankID].InterestIncome[currency], currency)
	commandlog.State("%s equity[%s]: %d %s", bankID, currency, w.Banks[bankID].Equity[currency], currency)
	return nil
}
