package sepa

import (
	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newRejectSepaCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reject-sepa payment_id reason",
		Short: "Reject a non-settled SEPA payment and refund the sender.",
		Args:  cobra.ExactArgs(2),
		RunE:  runRejectSepa,
	}
}

func runRejectSepa(cmd *cobra.Command, args []string) error {
	paymentID := args[0]
	reason := args[1]

	w, err := world.Load()
	if err != nil {
		return err
	}

	result, err := domain.RejectSepa(w, paymentID, reason)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "reject-sepa", args); err != nil {
		return err
	}

	payment := w.PaymentInstructions[paymentID]
	senderAccount := w.Accounts[result.SenderAccountID]
	commandlog.Action("Rejected SEPA payment: %s", paymentID)
	commandlog.State("Reason: %s", reason)
	commandlog.State("Status: %s", payment.Status)
	commandlog.State("%s booked_balance: %d %s", senderAccount.ID, senderAccount.BookedBalance, senderAccount.Currency)
	return nil
}
