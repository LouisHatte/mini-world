package payments

import (
	"strings"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newInternalTransferCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "internal-transfer sender_human_id recipient_human_id bank_id currency amount",
		Short: "Transfer deposits between two humans inside the same commercial bank.",
		Args:  cobra.ExactArgs(5),
		RunE:  runInternalTransfer,
	}
}

func runInternalTransfer(cmd *cobra.Command, args []string) error {
	senderHumanID := args[0]
	recipientHumanID := args[1]
	bankID := args[2]
	currency := strings.ToUpper(args[3])
	amount, err := parseAmount(args[4])
	if err != nil {
		return err
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	result, err := domain.InternalTransfer(w, senderHumanID, recipientHumanID, bankID, currency, amount)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "internal-transfer", args); err != nil {
		return err
	}

	commandlog.Action("Transferred %d %s from %s to %s inside %s", amount, currency, senderHumanID, recipientHumanID, bankID)
	logPaymentResult(w, result)
	return nil
}
