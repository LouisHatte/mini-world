package swift

import (
	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newSettleSwiftCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "settle-swift payment_id",
		Short: "Settle SWIFT payment through correspondent accounts.",
		Args:  cobra.ExactArgs(1),
		RunE:  runSettleSwift,
	}
}

func runSettleSwift(cmd *cobra.Command, args []string) error {
	paymentID := args[0]

	w, err := world.Load()
	if err != nil {
		return err
	}

	result, err := domain.SettleSwift(w, paymentID)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "settle-swift", args); err != nil {
		return err
	}

	payment := w.PaymentInstructions[paymentID]
	recipientAccount := w.Accounts[result.RecipientAccountID]
	commandlog.Action("Settled SWIFT payment: %s", paymentID)
	commandlog.State("Status: %s", payment.Status)
	logCorrespondentAccount(w, result.CorrespondentAccountID)
	commandlog.State("%s booked_balance: %d %s", recipientAccount.ID, recipientAccount.BookedBalance, recipientAccount.Currency)
	return nil
}
