package reserves

import (
	"fmt"
	"strconv"
	"strings"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newLendReservesCommand() *cobra.Command {
	var collateralAssetID string

	cmd := &cobra.Command{
		Use:   "lend-reserves central_bank_id bank_id currency amount --collateral asset_id",
		Short: "Central bank lends reserves against collateral.",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLendReserves(args, collateralAssetID)
		},
	}

	cmd.Flags().StringVar(&collateralAssetID, "collateral", "", "Collateral asset ID.")
	_ = cmd.MarkFlagRequired("collateral")

	return cmd
}

func runLendReserves(args []string, collateralAssetID string) error {
	centralBankID := args[0]
	bankID := args[1]
	currency := strings.ToUpper(args[2])
	amount, err := strconv.Atoi(args[3])
	if err != nil {
		return fmt.Errorf("amount must be an integer")
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	reserveLoanID, err := domain.LendReserves(w, centralBankID, bankID, currency, amount, collateralAssetID)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "lend-reserves", append(args, "--collateral", collateralAssetID)); err != nil {
		return err
	}

	centralBank := w.CentralBanks[centralBankID]
	bank := w.Banks[bankID]
	reserveLoan := w.ReserveLoans[reserveLoanID]
	commandlog.Action("Lent %d %s reserves from %s to %s", amount, currency, centralBankID, bankID)
	commandlog.State("Reserve loan: %s", reserveLoanID)
	commandlog.State("Collateral: %s", collateralAssetID)
	commandlog.State("%s reserve account for %s: %d %s", centralBankID, bankID, centralBank.ReserveAccounts[bankID], currency)
	commandlog.State("%s reserves at %s: %d %s", bankID, centralBankID, bank.ReserveBalances[centralBankID], currency)
	commandlog.State("%s loan to %s: %d %s", centralBankID, bankID, centralBank.LoansToBanks[bankID], currency)
	commandlog.State("%s loan from %s: %d %s", bankID, centralBankID, bank.LoansFromCentralBanks[centralBankID], currency)
	commandlog.State("%s outstanding: %d %s", reserveLoanID, reserveLoan.Outstanding, currency)
	return nil
}
