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

func newReserveTransferCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reserve-transfer central_bank_id from_bank_id to_bank_id currency amount",
		Short: "Transfer reserves between commercial banks.",
		Args:  cobra.ExactArgs(5),
		RunE:  runReserveTransfer,
	}
}

func runReserveTransfer(cmd *cobra.Command, args []string) error {
	centralBankID := args[0]
	fromBankID := args[1]
	toBankID := args[2]
	currency := strings.ToUpper(args[3])
	amount, err := strconv.Atoi(args[4])
	if err != nil {
		return fmt.Errorf("amount must be an integer")
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	if err := domain.ReserveTransfer(w, centralBankID, fromBankID, toBankID, currency, amount); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "reserve-transfer", args); err != nil {
		return err
	}

	centralBank := w.CentralBanks[centralBankID]
	fromBank := w.Banks[fromBankID]
	toBank := w.Banks[toBankID]
	commandlog.Action("Transferred %d %s reserves from %s to %s at %s", amount, currency, fromBankID, toBankID, centralBankID)
	commandlog.State("%s reserve account for %s: %d %s", centralBankID, fromBankID, centralBank.ReserveAccounts[fromBankID], currency)
	commandlog.State("%s reserve account for %s: %d %s", centralBankID, toBankID, centralBank.ReserveAccounts[toBankID], currency)
	commandlog.State("%s reserves at %s: %d %s", fromBankID, centralBankID, fromBank.ReserveBalances[centralBankID], currency)
	commandlog.State("%s reserves at %s: %d %s", toBankID, centralBankID, toBank.ReserveBalances[centralBankID], currency)
	return nil
}
