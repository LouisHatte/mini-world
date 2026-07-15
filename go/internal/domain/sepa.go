package domain

import (
	"fmt"
	"strings"

	"mini-world-go/internal/world"
)

type SepaCreateResult struct {
	PaymentID          string
	MessageID          string
	SenderAccountID    string
	RecipientAccountID string
}

type SepaSettlementResult struct {
	PaymentID        string
	SettlementID     string
	CentralBankID    string
	SettlementFailed bool
	MissingReserves  int
}

type SepaRejectResult struct {
	PaymentID       string
	SenderAccountID string
}

func SepaCreditTransfer(w *world.World, senderHumanID string, senderBankID string, recipientHumanID string, recipientBankID string, currency string, amount int) (SepaCreateResult, error) {
	return createSepaPayment(w, world.PaymentSepaCreditTransfer, senderHumanID, senderBankID, recipientHumanID, recipientBankID, currency, amount)
}

func SepaInstant(w *world.World, senderHumanID string, senderBankID string, recipientHumanID string, recipientBankID string, currency string, amount int) (SepaCreateResult, SepaSettlementResult, error) {
	createResult, err := createSepaPayment(w, world.PaymentSepaInstant, senderHumanID, senderBankID, recipientHumanID, recipientBankID, currency, amount)
	if err != nil {
		return SepaCreateResult{}, SepaSettlementResult{}, err
	}

	settlementResult, err := SettleSepa(w, createResult.PaymentID)
	if err != nil {
		return SepaCreateResult{}, SepaSettlementResult{}, err
	}

	return createResult, settlementResult, nil
}

func SettleSepa(w *world.World, paymentID string) (SepaSettlementResult, error) {
	payment, ok := w.PaymentInstructions[paymentID]
	if !ok {
		return SepaSettlementResult{}, fmt.Errorf("payment does not exist: %s", paymentID)
	}

	if payment.Rail != "SEPA" {
		return SepaSettlementResult{}, fmt.Errorf("payment is not a SEPA payment: %s", paymentID)
	}

	if payment.Status != world.PaymentInitiated {
		return SepaSettlementResult{}, fmt.Errorf("payment must be INITIATED, current status: %s", payment.Status)
	}

	senderBank, senderBankExists := w.Banks[payment.SenderBankID]
	if !senderBankExists {
		return SepaSettlementResult{}, fmt.Errorf("sender bank does not exist: %s", payment.SenderBankID)
	}

	recipientBank, recipientBankExists := w.Banks[payment.RecipientBankID]
	if !recipientBankExists {
		return SepaSettlementResult{}, fmt.Errorf("recipient bank does not exist: %s", payment.RecipientBankID)
	}

	recipientAccount, recipientAccountExists := w.Accounts[payment.RecipientAccountID]
	if !recipientAccountExists || recipientAccount.Status != world.AccountActive {
		return SepaSettlementResult{}, fmt.Errorf("recipient account is not active: %s", payment.RecipientAccountID)
	}

	centralBankID, centralBank, err := settlementCentralBank(w, payment.SenderBankID, payment.RecipientBankID, payment.Currency)
	if err != nil {
		return SepaSettlementResult{}, err
	}

	senderReserves, err := reserveBalance(centralBank, senderBank, centralBankID, payment.SenderBankID)
	if err != nil {
		return SepaSettlementResult{}, err
	}

	if _, err := reserveBalance(centralBank, recipientBank, centralBankID, payment.RecipientBankID); err != nil {
		return SepaSettlementResult{}, err
	}

	if senderReserves < payment.Amount {
		payment.Status = world.PaymentSettlementFailed
		return SepaSettlementResult{
			PaymentID:        paymentID,
			CentralBankID:    centralBankID,
			SettlementFailed: true,
			MissingReserves:  payment.Amount - senderReserves,
		}, nil
	}

	centralBank.ReserveAccounts[payment.SenderBankID] -= payment.Amount
	senderBank.ReserveBalances[centralBankID] -= payment.Amount
	centralBank.ReserveAccounts[payment.RecipientBankID] += payment.Amount
	recipientBank.ReserveBalances[centralBankID] += payment.Amount
	recipientAccount.BookedBalance += payment.Amount

	settlementID := nextSettlementID(w)
	w.Settlements[settlementID] = map[string]any{
		"id":              settlementID,
		"payment_id":      paymentID,
		"rail":            "T2",
		"central_bank_id": centralBankID,
		"status":          "SETTLED",
		"currency":        payment.Currency,
		"amount":          payment.Amount,
	}

	payment.CentralBankID = centralBankID
	payment.SettlementID = settlementID
	payment.Status = world.PaymentSettled

	return SepaSettlementResult{
		PaymentID:     paymentID,
		SettlementID:  settlementID,
		CentralBankID: centralBankID,
	}, nil
}

func RejectSepa(w *world.World, paymentID string, reason string) (SepaRejectResult, error) {
	payment, ok := w.PaymentInstructions[paymentID]
	if !ok {
		return SepaRejectResult{}, fmt.Errorf("payment does not exist: %s", paymentID)
	}

	if payment.Rail != "SEPA" {
		return SepaRejectResult{}, fmt.Errorf("payment is not a SEPA payment: %s", paymentID)
	}

	if payment.Status == world.PaymentSettled {
		return SepaRejectResult{}, fmt.Errorf("settled SEPA payment cannot be rejected: %s", paymentID)
	}

	if payment.Status == world.PaymentRejected {
		return SepaRejectResult{}, fmt.Errorf("SEPA payment is already rejected: %s", paymentID)
	}

	senderAccount, ok := w.Accounts[payment.SenderAccountID]
	if ok {
		senderAccount.BookedBalance += payment.Amount
	}

	payment.Status = world.PaymentRejected
	payment.ReturnReason = reason

	return SepaRejectResult{
		PaymentID:       paymentID,
		SenderAccountID: payment.SenderAccountID,
	}, nil
}

func createSepaPayment(w *world.World, paymentType world.PaymentType, senderHumanID string, senderBankID string, recipientHumanID string, recipientBankID string, currency string, amount int) (SepaCreateResult, error) {
	currency = strings.ToUpper(currency)
	if err := validatePaymentBasics(w, senderHumanID, recipientHumanID, currency, amount); err != nil {
		return SepaCreateResult{}, err
	}

	if senderBankID == recipientBankID {
		return SepaCreateResult{}, fmt.Errorf("use internal-transfer for same-bank payments")
	}

	if _, exists := w.Banks[senderBankID]; !exists {
		return SepaCreateResult{}, fmt.Errorf("sender bank does not exist: %s", senderBankID)
	}

	if _, exists := w.Banks[recipientBankID]; !exists {
		return SepaCreateResult{}, fmt.Errorf("recipient bank does not exist: %s", recipientBankID)
	}

	senderAccountID, senderAccount, err := activeAccount(w, senderHumanID, senderBankID, currency)
	if err != nil {
		return SepaCreateResult{}, err
	}

	recipientAccountID, _, err := activeAccount(w, recipientHumanID, recipientBankID, currency)
	if err != nil {
		return SepaCreateResult{}, err
	}

	if senderAccount.BookedBalance < amount {
		return SepaCreateResult{}, fmt.Errorf(
			"not enough deposits in %s. Available: %d %s",
			senderAccountID,
			senderAccount.BookedBalance,
			currency,
		)
	}

	senderAccount.BookedBalance -= amount

	paymentID := nextPaymentID(w)
	messageID := nextMessageID(w)
	w.PaymentInstructions[paymentID] = world.NewPaymentInstruction(
		paymentID,
		paymentType,
		senderHumanID,
		senderBankID,
		recipientHumanID,
		recipientBankID,
		senderAccountID,
		recipientAccountID,
		"",
		currency,
		amount,
	)

	payment := w.PaymentInstructions[paymentID]
	payment.Rail = "SEPA"
	payment.Status = world.PaymentInitiated
	payment.MessageID = messageID

	w.Messages[messageID] = map[string]any{
		"id":                messageID,
		"payment_id":        paymentID,
		"rail":              "SEPA",
		"sender_bank_id":    senderBankID,
		"recipient_bank_id": recipientBankID,
		"status":            "SENT",
	}

	return SepaCreateResult{
		PaymentID:          paymentID,
		MessageID:          messageID,
		SenderAccountID:    senderAccountID,
		RecipientAccountID: recipientAccountID,
	}, nil
}

func nextMessageID(w *world.World) string {
	for index := 1; ; index++ {
		messageID := fmt.Sprintf("msg_%06d", index)
		if _, exists := w.Messages[messageID]; !exists {
			return messageID
		}
	}
}

func nextSettlementID(w *world.World) string {
	for index := 1; ; index++ {
		settlementID := fmt.Sprintf("settlement_%06d", index)
		if _, exists := w.Settlements[settlementID]; !exists {
			return settlementID
		}
	}
}
