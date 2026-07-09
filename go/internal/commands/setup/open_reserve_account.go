package setup

import (
	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newOpenReserveAccountCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "open-reserve-account central_bank_id bank_id",
		Short: "Open a commercial bank reserve account at a central bank.",
		Args:  cobra.ExactArgs(2),
		RunE:  runOpenReserveAccount,
	}
}

func runOpenReserveAccount(cmd *cobra.Command, args []string) error {
	centralBankID := args[0]
	bankID := args[1]

	w, err := world.Load()
	if err != nil {
		return err
	}

	if err := domain.OpenReserveAccount(w, centralBankID, bankID); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "open-reserve-account", args); err != nil {
		return err
	}

	centralBank := w.CentralBanks[centralBankID]
	bank := w.Banks[bankID]
	currency := centralBank.Currency
	commandlog.Action("Opened reserve account: %s at %s", bankID, centralBankID)
	commandlog.State("%s reserve account for %s: %d %s", centralBankID, bankID, centralBank.ReserveAccounts[bankID], currency)
	commandlog.State("%s reserves at %s: %d %s", bankID, centralBankID, bank.ReserveBalances[centralBankID], currency)
	return nil
}
