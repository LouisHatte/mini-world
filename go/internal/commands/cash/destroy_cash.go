package cash

import (
	"fmt"
	"strconv"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newDestroyCashCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "destroy-cash central_bank_id amount",
		Short: "Destroy physical cash held in a central bank vault.",
		Args:  cobra.ExactArgs(2),
		RunE:  runDestroyCash,
	}
}

func runDestroyCash(cmd *cobra.Command, args []string) error {
	centralBankID := args[0]
	amount, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("amount must be an integer")
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	if err := domain.DestroyCash(w, centralBankID, amount); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "destroy-cash", args); err != nil {
		return err
	}

	centralBank := w.CentralBanks[centralBankID]
	currency := centralBank.Currency
	commandlog.Action("Destroyed %d %s cash at %s", amount, currency, centralBankID)
	commandlog.State("%s cash_issued: %d %s", centralBankID, centralBank.CashIssued, currency)
	commandlog.State("%s cash_vault: %d %s", centralBankID, centralBank.CashVault, currency)
	return nil
}
