package swift

import (
	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newRejectSwiftCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reject-swift payment_id reason",
		Short: "Reject a non-settled SWIFT payment.",
		Args:  cobra.ExactArgs(2),
		RunE:  runRejectSwift,
	}
}

func runRejectSwift(cmd *cobra.Command, args []string) error {
	paymentID := args[0]
	reason := args[1]

	w, err := world.Load()
	if err != nil {
		return err
	}

	result, err := domain.RejectSwift(w, paymentID, reason)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "reject-swift", args); err != nil {
		return err
	}

	payment := w.PaymentInstructions[paymentID]
	senderAccount := w.Accounts[result.SenderAccountID]
	commandlog.Action("Rejected SWIFT payment: %s", paymentID)
	commandlog.State("Reason: %s", reason)
	commandlog.State("Status: %s", payment.Status)
	commandlog.State("%s booked_balance: %d %s", senderAccount.ID, senderAccount.BookedBalance, senderAccount.Currency)
	return nil
}
