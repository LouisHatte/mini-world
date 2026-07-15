package domain

import (
	"fmt"
	"strings"

	"mini-world-go/internal/world"
)

type PaymentResult struct {
	PaymentID          string
	Type               world.PaymentType
	SenderAccountID    string
	RecipientAccountID string
	CentralBankID      string
}

func InternalTransfer(w *world.World, senderHumanID string, recipientHumanID string, bankID string, currency string, amount int) (PaymentResult, error) {
	currency = strings.ToUpper(currency)
	if err := validatePaymentBasics(w, senderHumanID, recipientHumanID, currency, amount); err != nil {
		return PaymentResult{}, err
	}

	if _, exists := w.Banks[bankID]; !exists {
		return PaymentResult{}, fmt.Errorf("bank does not exist: %s", bankID)
	}

	senderAccountID, senderAccount, err := activeAccount(w, senderHumanID, bankID, currency)
	if err != nil {
		return PaymentResult{}, err
	}

	recipientAccountID, recipientAccount, err := activeAccount(w, recipientHumanID, bankID, currency)
	if err != nil {
		return PaymentResult{}, err
	}

	if senderAccount.BookedBalance < amount {
		return PaymentResult{}, fmt.Errorf(
			"not enough deposits in %s. Available: %d %s",
			senderAccountID,
			senderAccount.BookedBalance,
			currency,
		)
	}

	senderAccount.BookedBalance -= amount
	recipientAccount.BookedBalance += amount

	paymentID := nextPaymentID(w)
	w.PaymentInstructions[paymentID] = world.NewPaymentInstruction(
		paymentID,
		world.PaymentInternal,
		senderHumanID,
		bankID,
		recipientHumanID,
		bankID,
		senderAccountID,
		recipientAccountID,
		"",
		currency,
		amount,
	)

	return PaymentResult{
		PaymentID:          paymentID,
		Type:               world.PaymentInternal,
		SenderAccountID:    senderAccountID,
		RecipientAccountID: recipientAccountID,
	}, nil
}

func InterbankPayment(w *world.World, senderHumanID string, senderBankID string, recipientHumanID string, recipientBankID string, currency string, amount int) (PaymentResult, error) {
	currency = strings.ToUpper(currency)
	if err := validatePaymentBasics(w, senderHumanID, recipientHumanID, currency, amount); err != nil {
		return PaymentResult{}, err
	}

	if senderBankID == recipientBankID {
		return PaymentResult{}, fmt.Errorf("sender bank and recipient bank must be different")
	}

	senderBank, senderBankExists := w.Banks[senderBankID]
	if !senderBankExists {
		return PaymentResult{}, fmt.Errorf("sender bank does not exist: %s", senderBankID)
	}

	recipientBank, recipientBankExists := w.Banks[recipientBankID]
	if !recipientBankExists {
		return PaymentResult{}, fmt.Errorf("recipient bank does not exist: %s", recipientBankID)
	}

	senderAccountID, senderAccount, err := activeAccount(w, senderHumanID, senderBankID, currency)
	if err != nil {
		return PaymentResult{}, err
	}

	recipientAccountID, recipientAccount, err := activeAccount(w, recipientHumanID, recipientBankID, currency)
	if err != nil {
		return PaymentResult{}, err
	}

	if senderAccount.BookedBalance < amount {
		return PaymentResult{}, fmt.Errorf(
			"not enough deposits in %s. Available: %d %s",
			senderAccountID,
			senderAccount.BookedBalance,
			currency,
		)
	}

	centralBankID, centralBank, err := settlementCentralBank(w, senderBankID, recipientBankID, currency)
	if err != nil {
		return PaymentResult{}, err
	}

	senderReserves, err := reserveBalance(centralBank, senderBank, centralBankID, senderBankID)
	if err != nil {
		return PaymentResult{}, err
	}

	if _, err := reserveBalance(centralBank, recipientBank, centralBankID, recipientBankID); err != nil {
		return PaymentResult{}, err
	}

	if senderReserves < amount {
		return PaymentResult{}, fmt.Errorf(
			"not enough reserves for %s at %s. Available: %d %s",
			senderBankID,
			centralBankID,
			senderReserves,
			currency,
		)
	}

	senderAccount.BookedBalance -= amount
	recipientAccount.BookedBalance += amount
	centralBank.ReserveAccounts[senderBankID] -= amount
	senderBank.ReserveBalances[centralBankID] -= amount
	centralBank.ReserveAccounts[recipientBankID] += amount
	recipientBank.ReserveBalances[centralBankID] += amount

	paymentID := nextPaymentID(w)
	w.PaymentInstructions[paymentID] = world.NewPaymentInstruction(
		paymentID,
		world.PaymentInterbank,
		senderHumanID,
		senderBankID,
		recipientHumanID,
		recipientBankID,
		senderAccountID,
		recipientAccountID,
		centralBankID,
		currency,
		amount,
	)

	return PaymentResult{
		PaymentID:          paymentID,
		Type:               world.PaymentInterbank,
		SenderAccountID:    senderAccountID,
		RecipientAccountID: recipientAccountID,
		CentralBankID:      centralBankID,
	}, nil
}

func Pay(w *world.World, senderHumanID string, senderBankID string, recipientHumanID string, recipientBankID string, currency string, amount int) (PaymentResult, error) {
	if senderBankID == recipientBankID {
		return InternalTransfer(w, senderHumanID, recipientHumanID, senderBankID, currency, amount)
	}

	return InterbankPayment(w, senderHumanID, senderBankID, recipientHumanID, recipientBankID, currency, amount)
}

func validatePaymentBasics(w *world.World, senderHumanID string, recipientHumanID string, currency string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if senderHumanID == recipientHumanID {
		return fmt.Errorf("sender and recipient must be different humans")
	}

	if _, exists := w.Humans[senderHumanID]; !exists {
		return fmt.Errorf("sender human does not exist: %s", senderHumanID)
	}

	if _, exists := w.Humans[recipientHumanID]; !exists {
		return fmt.Errorf("recipient human does not exist: %s", recipientHumanID)
	}

	if _, exists := w.Currencies[currency]; !exists {
		return fmt.Errorf("currency does not exist: %s", currency)
	}

	return nil
}

func activeAccount(w *world.World, humanID string, bankID string, currency string) (string, *world.Account, error) {
	accountID := BuildAccountID(bankID, humanID, currency)
	account, ok := w.Accounts[accountID]

	if !ok || account.Status != world.AccountActive {
		return "", nil, fmt.Errorf("no active %s account for %s at %s", currency, humanID, bankID)
	}

	return accountID, account, nil
}

func settlementCentralBank(w *world.World, senderBankID string, recipientBankID string, currency string) (string, *world.CentralBank, error) {
	matchingCentralBankID := ""
	var matchingCentralBank *world.CentralBank

	for centralBankID, centralBank := range w.CentralBanks {
		if centralBank.Currency != currency {
			continue
		}

		if _, senderHasReserveAccount := centralBank.ReserveAccounts[senderBankID]; !senderHasReserveAccount {
			continue
		}

		if _, recipientHasReserveAccount := centralBank.ReserveAccounts[recipientBankID]; !recipientHasReserveAccount {
			continue
		}

		if matchingCentralBankID != "" {
			return "", nil, fmt.Errorf("multiple settlement central banks found for %s between %s and %s", currency, senderBankID, recipientBankID)
		}

		matchingCentralBankID = centralBankID
		matchingCentralBank = centralBank
	}

	if matchingCentralBank == nil {
		return "", nil, fmt.Errorf("no settlement central bank found for %s between %s and %s", currency, senderBankID, recipientBankID)
	}

	return matchingCentralBankID, matchingCentralBank, nil
}

func nextPaymentID(w *world.World) string {
	for index := 1; ; index++ {
		paymentID := fmt.Sprintf("payment_%06d", index)
		if _, exists := w.PaymentInstructions[paymentID]; !exists {
			return paymentID
		}
	}
}
