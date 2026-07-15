package sepa

import (
	"strings"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newSepaCreditTransferCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sepa-credit-transfer sender_human_id sender_bank_id recipient_human_id recipient_bank_id currency amount",
		Short: "Create a SEPA credit transfer instruction.",
		Args:  cobra.ExactArgs(6),
		RunE:  runSepaCreditTransfer,
	}
}

func runSepaCreditTransfer(cmd *cobra.Command, args []string) error {
	senderHumanID := args[0]
	senderBankID := args[1]
	recipientHumanID := args[2]
	recipientBankID := args[3]
	currency := strings.ToUpper(args[4])
	amount, err := parseAmount(args[5])
	if err != nil {
		return err
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	result, err := domain.SepaCreditTransfer(w, senderHumanID, senderBankID, recipientHumanID, recipientBankID, currency, amount)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "sepa-credit-transfer", args); err != nil {
		return err
	}

	commandlog.Action("Created SEPA credit transfer of %d %s from %s at %s to %s at %s", amount, currency, senderHumanID, senderBankID, recipientHumanID, recipientBankID)
	logSepaCreateResult(w, result)
	return nil
}
