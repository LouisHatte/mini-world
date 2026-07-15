package payments

import (
	"strings"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newPayCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "pay sender_human_id sender_bank_id recipient_human_id recipient_bank_id currency amount",
		Short: "High-level payment command, auto-detects internal or interbank payment.",
		Args:  cobra.ExactArgs(6),
		RunE:  runPay,
	}
}

func runPay(cmd *cobra.Command, args []string) error {
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

	result, err := domain.Pay(w, senderHumanID, senderBankID, recipientHumanID, recipientBankID, currency, amount)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "pay", args); err != nil {
		return err
	}

	commandlog.Action("Paid %d %s from %s at %s to %s at %s", amount, currency, senderHumanID, senderBankID, recipientHumanID, recipientBankID)
	logPaymentResult(w, result)
	return nil
}
