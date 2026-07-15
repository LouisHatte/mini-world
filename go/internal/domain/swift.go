package domain

import (
	"fmt"
	"strings"

	"mini-world-go/internal/world"
)

type CorrespondentFundingResult struct {
	CorrespondentAccountID string
	CentralBankID          string
}

type SwiftCreateResult struct {
	PaymentID          string
	MessageID          string
	SenderAccountID    string
	RecipientAccountID string
}

type SwiftSettlementResult struct {
	PaymentID              string
	CorrespondentAccountID string
	RecipientAccountID     string
}

type SwiftRejectResult struct {
	PaymentID       string
	SenderAccountID string
}

func OpenCorrespondentAccount(w *world.World, ownerBankID string, correspondentBankID string, currency string) (string, error) {
	currency = strings.ToUpper(currency)

	if ownerBankID == correspondentBankID {
		return "", fmt.Errorf("owner bank and correspondent bank must be different")
	}

	ownerBank, exists := w.Banks[ownerBankID]
	if !exists {
		return "", fmt.Errorf("owner bank does not exist: %s", ownerBankID)
	}

	correspondentBank, exists := w.Banks[correspondentBankID]
	if !exists {
		return "", fmt.Errorf("correspondent bank does not exist: %s", correspondentBankID)
	}

	if _, exists := w.Currencies[currency]; !exists {
		return "", fmt.Errorf("currency does not exist: %s", currency)
	}

	accountID := CorrespondentAccountID(ownerBankID, correspondentBankID, currency)
	if _, exists := w.CorrespondentAccounts[accountID]; exists {
		return "", fmt.Errorf("correspondent account already exists: %s", accountID)
	}

	w.CorrespondentAccounts[accountID] = world.NewCorrespondentAccount(accountID, ownerBankID, correspondentBankID, currency)
	ownerBank.NostroAccounts = append(ownerBank.NostroAccounts, accountID)
	correspondentBank.VostroAccounts = append(correspondentBank.VostroAccounts, accountID)

	return accountID, nil
}

func FundCorrespondentAccount(w *world.World, correspondentAccountID string, amount int) (CorrespondentFundingResult, error) {
	if amount <= 0 {
		return CorrespondentFundingResult{}, fmt.Errorf("amount must be greater than 0")
	}

	correspondentAccount, ok := w.CorrespondentAccounts[correspondentAccountID]
	if !ok {
		return CorrespondentFundingResult{}, fmt.Errorf("correspondent account does not exist: %s", correspondentAccountID)
	}

	if correspondentAccount.Status != world.CorrespondentAccountActive {
		return CorrespondentFundingResult{}, fmt.Errorf("correspondent account is not active: %s", correspondentAccountID)
	}

	ownerBank := w.Banks[correspondentAccount.OwnerBankID]
	correspondentBank := w.Banks[correspondentAccount.CorrespondentBankID]
	centralBankID, centralBank, err := centralBankForCurrency(w, correspondentAccount.Currency)
	if err != nil {
		return CorrespondentFundingResult{}, err
	}

	ownerReserves, err := reserveBalance(centralBank, ownerBank, centralBankID, correspondentAccount.OwnerBankID)
	if err != nil {
		return CorrespondentFundingResult{}, err
	}

	if _, err := reserveBalance(centralBank, correspondentBank, centralBankID, correspondentAccount.CorrespondentBankID); err != nil {
		return CorrespondentFundingResult{}, err
	}

	if ownerReserves < amount {
		return CorrespondentFundingResult{}, fmt.Errorf(
			"not enough reserves for %s at %s. Available: %d %s",
			correspondentAccount.OwnerBankID,
			centralBankID,
			ownerReserves,
			correspondentAccount.Currency,
		)
	}

	centralBank.ReserveAccounts[correspondentAccount.OwnerBankID] -= amount
	ownerBank.ReserveBalances[centralBankID] -= amount
	centralBank.ReserveAccounts[correspondentAccount.CorrespondentBankID] += amount
	correspondentBank.ReserveBalances[centralBankID] += amount
	correspondentAccount.NostroBalance += amount
	correspondentAccount.VostroBalance += amount

	return CorrespondentFundingResult{
		CorrespondentAccountID: correspondentAccountID,
		CentralBankID:          centralBankID,
	}, nil
}

func SwiftMT103(w *world.World, senderHumanID string, senderBankID string, recipientHumanID string, recipientBankID string, currency string, amount int) (SwiftCreateResult, error) {
	currency = strings.ToUpper(currency)
	if err := validatePaymentBasics(w, senderHumanID, recipientHumanID, currency, amount); err != nil {
		return SwiftCreateResult{}, err
	}

	if senderBankID == recipientBankID {
		return SwiftCreateResult{}, fmt.Errorf("use internal-transfer for same-bank payments")
	}

	if _, exists := w.Banks[senderBankID]; !exists {
		return SwiftCreateResult{}, fmt.Errorf("sender bank does not exist: %s", senderBankID)
	}

	if _, exists := w.Banks[recipientBankID]; !exists {
		return SwiftCreateResult{}, fmt.Errorf("recipient bank does not exist: %s", recipientBankID)
	}

	senderAccountID, senderAccount, err := activeAccount(w, senderHumanID, senderBankID, currency)
	if err != nil {
		return SwiftCreateResult{}, err
	}

	recipientAccountID, _, err := activeAccount(w, recipientHumanID, recipientBankID, currency)
	if err != nil {
		return SwiftCreateResult{}, err
	}

	if senderAccount.BookedBalance < amount {
		return SwiftCreateResult{}, fmt.Errorf(
			"not enough deposits in %s. Available: %d %s",
			senderAccountID,
			senderAccount.BookedBalance,
			currency,
		)
	}

	correspondentAccountID := CorrespondentAccountID(senderBankID, recipientBankID, currency)
	if _, exists := w.CorrespondentAccounts[correspondentAccountID]; !exists {
		return SwiftCreateResult{}, fmt.Errorf("correspondent account does not exist: %s", correspondentAccountID)
	}

	senderAccount.BookedBalance -= amount

	paymentID := nextPaymentID(w)
	messageID := nextMessageID(w)
	w.PaymentInstructions[paymentID] = world.NewPaymentInstruction(
		paymentID,
		world.PaymentSwiftMT103,
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
	payment.Rail = "SWIFT"
	payment.Status = world.PaymentInstructed
	payment.MessageID = messageID

	w.Messages[messageID] = map[string]any{
		"id":                messageID,
		"payment_id":        paymentID,
		"rail":              "SWIFT",
		"type":              "MT103",
		"sender_bank_id":    senderBankID,
		"recipient_bank_id": recipientBankID,
		"status":            "SENT",
	}

	return SwiftCreateResult{
		PaymentID:          paymentID,
		MessageID:          messageID,
		SenderAccountID:    senderAccountID,
		RecipientAccountID: recipientAccountID,
	}, nil
}

func SettleSwift(w *world.World, paymentID string) (SwiftSettlementResult, error) {
	payment, ok := w.PaymentInstructions[paymentID]
	if !ok {
		return SwiftSettlementResult{}, fmt.Errorf("payment does not exist: %s", paymentID)
	}

	if payment.Rail != "SWIFT" {
		return SwiftSettlementResult{}, fmt.Errorf("payment is not a SWIFT payment: %s", paymentID)
	}

	if payment.Status != world.PaymentInstructed {
		return SwiftSettlementResult{}, fmt.Errorf("payment must be INSTRUCTED, current status: %s", payment.Status)
	}

	correspondentAccountID := CorrespondentAccountID(payment.SenderBankID, payment.RecipientBankID, payment.Currency)
	correspondentAccount, ok := w.CorrespondentAccounts[correspondentAccountID]
	if !ok {
		return SwiftSettlementResult{}, fmt.Errorf("correspondent account does not exist: %s", correspondentAccountID)
	}

	if correspondentAccount.NostroBalance < payment.Amount || correspondentAccount.VostroBalance < payment.Amount {
		return SwiftSettlementResult{}, fmt.Errorf(
			"not enough correspondent balance in %s. Available: %d %s",
			correspondentAccountID,
			correspondentAccount.NostroBalance,
			payment.Currency,
		)
	}

	recipientAccount, ok := w.Accounts[payment.RecipientAccountID]
	if !ok || recipientAccount.Status != world.AccountActive {
		return SwiftSettlementResult{}, fmt.Errorf("recipient account is not active: %s", payment.RecipientAccountID)
	}

	correspondentAccount.NostroBalance -= payment.Amount
	correspondentAccount.VostroBalance -= payment.Amount
	recipientAccount.BookedBalance += payment.Amount
	payment.Status = world.PaymentSettled

	return SwiftSettlementResult{
		PaymentID:              paymentID,
		CorrespondentAccountID: correspondentAccountID,
		RecipientAccountID:     payment.RecipientAccountID,
	}, nil
}

func RejectSwift(w *world.World, paymentID string, reason string) (SwiftRejectResult, error) {
	payment, ok := w.PaymentInstructions[paymentID]
	if !ok {
		return SwiftRejectResult{}, fmt.Errorf("payment does not exist: %s", paymentID)
	}

	if payment.Rail != "SWIFT" {
		return SwiftRejectResult{}, fmt.Errorf("payment is not a SWIFT payment: %s", paymentID)
	}

	if payment.Status == world.PaymentSettled {
		return SwiftRejectResult{}, fmt.Errorf("settled SWIFT payment cannot be rejected: %s", paymentID)
	}

	if payment.Status == world.PaymentRejected {
		return SwiftRejectResult{}, fmt.Errorf("SWIFT payment is already rejected: %s", paymentID)
	}

	senderAccount, ok := w.Accounts[payment.SenderAccountID]
	if ok {
		senderAccount.BookedBalance += payment.Amount
	}

	payment.Status = world.PaymentRejected
	payment.ReturnReason = reason

	return SwiftRejectResult{
		PaymentID:       paymentID,
		SenderAccountID: payment.SenderAccountID,
	}, nil
}

func CorrespondentAccountID(ownerBankID string, correspondentBankID string, currency string) string {
	return fmt.Sprintf("corr_%s_%s_%s", ownerBankID, correspondentBankID, strings.ToLower(currency))
}
