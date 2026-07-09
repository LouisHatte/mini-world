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

func newReturnCashCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "return-cash central_bank_id bank_id amount",
		Short: "Return physical cash from a commercial bank to a central bank in exchange for reserves.",
		Args:  cobra.ExactArgs(3),
		RunE:  runReturnCash,
	}
}

func runReturnCash(cmd *cobra.Command, args []string) error {
	centralBankID := args[0]
	bankID := args[1]
	amount, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("amount must be an integer")
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	if err := domain.ReturnCash(w, centralBankID, bankID, amount); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "return-cash", args); err != nil {
		return err
	}

	centralBank := w.CentralBanks[centralBankID]
	bank := w.Banks[bankID]
	currency := centralBank.Currency
	commandlog.Action("Returned %d %s cash from %s to %s", amount, currency, bankID, centralBankID)
	commandlog.State("%s cash_vault[%s]: %d %s", bankID, currency, bank.CashVault[currency], currency)
	commandlog.State("%s cash_vault: %d %s", centralBankID, centralBank.CashVault, currency)
	commandlog.State("%s reserve account for %s: %d %s", centralBankID, bankID, centralBank.ReserveAccounts[bankID], currency)
	commandlog.State("%s reserves at %s: %d %s", bankID, centralBankID, bank.ReserveBalances[centralBankID], currency)
	return nil
}
