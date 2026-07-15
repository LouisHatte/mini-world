package loans

import (
	"strings"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newGrantLoanCommand() *cobra.Command {
	var collateralAssetID string

	cmd := &cobra.Command{
		Use:   "grant-loan bank_id human_id currency amount",
		Short: "Grant a commercial bank loan and credit the borrower's deposit account.",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGrantLoan(args, collateralAssetID)
		},
	}

	cmd.Flags().StringVar(&collateralAssetID, "collateral", "", "Collateral asset ID.")

	return cmd
}

func runGrantLoan(args []string, collateralAssetID string) error {
	bankID := args[0]
	humanID := args[1]
	currency := strings.ToUpper(args[2])
	amount, err := parseAmount(args[3])
	if err != nil {
		return err
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	loanID, accountID, err := domain.GrantLoan(w, bankID, humanID, currency, amount, collateralAssetID)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	historyArgs := args
	if collateralAssetID != "" {
		historyArgs = append(historyArgs, "--collateral", collateralAssetID)
	}

	if err := commandrun.SaveWithHistory(w, "grant-loan", historyArgs); err != nil {
		return err
	}

	commandlog.Action("Granted %d %s loan from %s to %s", amount, currency, bankID, humanID)
	commandlog.State("Loan: %s", loanID)
	logLoanState(w, loanID)
	if collateralAssetID != "" {
		commandlog.State("Collateral: %s", collateralAssetID)
		commandlog.State("%s pledged_to_bank_id: %s", collateralAssetID, w.Assets[collateralAssetID].PledgedToBankID)
		commandlog.State("%s collateral_for_customer_loan_id: %s", collateralAssetID, w.Assets[collateralAssetID].CollateralForCustomerLoanID)
	}
	commandlog.State("%s booked_balance: %d %s", accountID, w.Accounts[accountID].BookedBalance, currency)
	return nil
}
