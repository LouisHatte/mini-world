package payments

import (
	"fmt"
	"strconv"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"
)

func parseAmount(value string) (int, error) {
	amount, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("amount must be an integer")
	}

	return amount, nil
}

func logPaymentResult(w *world.World, result domain.PaymentResult) {
	payment := w.PaymentInstructions[result.PaymentID]
	senderAccount := w.Accounts[result.SenderAccountID]
	recipientAccount := w.Accounts[result.RecipientAccountID]

	commandlog.State("Payment: %s", result.PaymentID)
	commandlog.State("Type: %s", result.Type)
	commandlog.State("%s booked_balance: %d %s", senderAccount.ID, senderAccount.BookedBalance, senderAccount.Currency)
	commandlog.State("%s booked_balance: %d %s", recipientAccount.ID, recipientAccount.BookedBalance, recipientAccount.Currency)

	if result.CentralBankID == "" {
		return
	}

	centralBank := w.CentralBanks[result.CentralBankID]
	senderBank := w.Banks[payment.SenderBankID]
	recipientBank := w.Banks[payment.RecipientBankID]

	commandlog.State("%s reserve account for %s: %d %s", result.CentralBankID, payment.SenderBankID, centralBank.ReserveAccounts[payment.SenderBankID], payment.Currency)
	commandlog.State("%s reserve account for %s: %d %s", result.CentralBankID, payment.RecipientBankID, centralBank.ReserveAccounts[payment.RecipientBankID], payment.Currency)
	commandlog.State("%s reserves at %s: %d %s", payment.SenderBankID, result.CentralBankID, senderBank.ReserveBalances[result.CentralBankID], payment.Currency)
	commandlog.State("%s reserves at %s: %d %s", payment.RecipientBankID, result.CentralBankID, recipientBank.ReserveBalances[result.CentralBankID], payment.Currency)
}
