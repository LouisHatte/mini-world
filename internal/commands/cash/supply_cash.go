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

func newSupplyCashCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "supply-cash central_bank_id bank_id amount",
		Short: "Move physical cash from a central bank vault to a commercial bank vault.",
		Args:  cobra.ExactArgs(3),
		RunE:  runSupplyCash,
	}
}

func runSupplyCash(cmd *cobra.Command, args []string) error {
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

	if err := domain.SupplyCash(w, centralBankID, bankID, amount); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "supply-cash", args); err != nil {
		return err
	}

	centralBank := w.CentralBanks[centralBankID]
	bank := w.Banks[bankID]
	currency := centralBank.Currency
	commandlog.Action("Supplied %d %s cash from %s to %s", amount, currency, centralBankID, bankID)
	commandlog.State("%s cash_vault: %d %s", centralBankID, centralBank.CashVault, currency)
	commandlog.State("%s reserve account for %s: %d %s", centralBankID, bankID, centralBank.ReserveAccounts[bankID], currency)
	commandlog.State("%s cash_vault[%s]: %d %s", bankID, currency, bank.CashVault[currency], currency)
	commandlog.State("%s reserves at %s: %d %s", bankID, centralBankID, bank.ReserveBalances[centralBankID], currency)
	return nil
}
