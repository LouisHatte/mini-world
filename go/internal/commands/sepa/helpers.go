package sepa

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

func logSepaCreateResult(w *world.World, result domain.SepaCreateResult) {
	payment := w.PaymentInstructions[result.PaymentID]
	senderAccount := w.Accounts[result.SenderAccountID]

	commandlog.State("Payment: %s", result.PaymentID)
	commandlog.State("Message: %s", result.MessageID)
	commandlog.State("Status: %s", payment.Status)
	commandlog.State("%s booked_balance: %d %s", senderAccount.ID, senderAccount.BookedBalance, senderAccount.Currency)
}

func logSepaSettlementResult(w *world.World, result domain.SepaSettlementResult) {
	payment := w.PaymentInstructions[result.PaymentID]

	if result.SettlementFailed {
		commandlog.State("Status: %s", payment.Status)
		commandlog.State("Missing reserves: %d %s", result.MissingReserves, payment.Currency)
		return
	}

	centralBank := w.CentralBanks[result.CentralBankID]
	senderBank := w.Banks[payment.SenderBankID]
	recipientBank := w.Banks[payment.RecipientBankID]
	recipientAccount := w.Accounts[payment.RecipientAccountID]

	commandlog.State("Settlement: %s", result.SettlementID)
	commandlog.State("Status: %s", payment.Status)
	commandlog.State("%s booked_balance: %d %s", recipientAccount.ID, recipientAccount.BookedBalance, payment.Currency)
	commandlog.State("%s reserve account for %s: %d %s", result.CentralBankID, payment.SenderBankID, centralBank.ReserveAccounts[payment.SenderBankID], payment.Currency)
	commandlog.State("%s reserve account for %s: %d %s", result.CentralBankID, payment.RecipientBankID, centralBank.ReserveAccounts[payment.RecipientBankID], payment.Currency)
	commandlog.State("%s reserves at %s: %d %s", payment.SenderBankID, result.CentralBankID, senderBank.ReserveBalances[result.CentralBankID], payment.Currency)
	commandlog.State("%s reserves at %s: %d %s", payment.RecipientBankID, result.CentralBankID, recipientBank.ReserveBalances[result.CentralBankID], payment.Currency)
}
