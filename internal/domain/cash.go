package domain

import (
	"fmt"
	"strings"

	"mini-world-go/internal/world"
)

func IssueCash(w *world.World, centralBankID string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	centralBank, ok := w.CentralBanks[centralBankID]
	if !ok {
		return fmt.Errorf("central bank does not exist: %s", centralBankID)
	}

	centralBank.CashIssued += amount
	centralBank.CashVault += amount

	return nil
}

func SeedCash(w *world.World, centralBankID string, humanID string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	centralBank, ok := w.CentralBanks[centralBankID]
	if !ok {
		return fmt.Errorf("central bank does not exist: %s", centralBankID)
	}

	human, ok := w.Humans[humanID]
	if !ok {
		return fmt.Errorf("human does not exist: %s", humanID)
	}

	currency := centralBank.Currency
	if centralBank.CashVault < amount {
		return fmt.Errorf(
			"not enough cash in %s vault. Available: %d %s",
			centralBankID,
			centralBank.CashVault,
			currency,
		)
	}

	centralBank.CashVault -= amount
	human.CashWallet[currency] += amount

	return nil
}

func TransferCash(w *world.World, sourceHumanID string, targetHumanID string, currency string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if sourceHumanID == targetHumanID {
		return fmt.Errorf("source human and target human must be different")
	}

	sourceHuman, ok := w.Humans[sourceHumanID]
	if !ok {
		return fmt.Errorf("source human does not exist: %s", sourceHumanID)
	}

	targetHuman, ok := w.Humans[targetHumanID]
	if !ok {
		return fmt.Errorf("target human does not exist: %s", targetHumanID)
	}

	currency = strings.ToUpper(currency)
	sourceCash := sourceHuman.CashWallet[currency]

	if sourceCash < amount {
		return fmt.Errorf(
			"not enough cash in %s's wallet. Available: %d %s",
			sourceHumanID,
			sourceCash,
			currency,
		)
	}

	sourceHuman.CashWallet[currency] = sourceCash - amount
	targetHuman.CashWallet[currency] += amount

	return nil
}

func DepositCash(w *world.World, humanID string, bankID string, currency string, amount int) (string, error) {
	if amount <= 0 {
		return "", fmt.Errorf("amount must be greater than 0")
	}

	human, ok := w.Humans[humanID]
	if !ok {
		return "", fmt.Errorf("human does not exist: %s", humanID)
	}

	bank, ok := w.Banks[bankID]
	if !ok {
		return "", fmt.Errorf("bank does not exist: %s", bankID)
	}

	currency = strings.ToUpper(currency)
	humanCash := human.CashWallet[currency]

	if humanCash < amount {
		return "", fmt.Errorf(
			"not enough cash in %s's wallet. Available: %d %s",
			humanID,
			humanCash,
			currency,
		)
	}

	accountID, ok := ActiveAccountID(w, humanID, bankID, currency)
	if !ok {
		return "", fmt.Errorf("no active %s account for %s at %s", currency, humanID, bankID)
	}

	account := w.Accounts[accountID]
	human.CashWallet[currency] = humanCash - amount
	bank.CashVault[currency] += amount
	account.BookedBalance += amount

	return accountID, nil
}

func WithdrawCash(w *world.World, humanID string, bankID string, currency string, amount int) (string, error) {
	if amount <= 0 {
		return "", fmt.Errorf("amount must be greater than 0")
	}

	if _, ok := w.Humans[humanID]; !ok {
		return "", fmt.Errorf("human does not exist: %s", humanID)
	}

	bank, ok := w.Banks[bankID]
	if !ok {
		return "", fmt.Errorf("bank does not exist: %s", bankID)
	}

	currency = strings.ToUpper(currency)
	accountID, ok := ActiveAccountID(w, humanID, bankID, currency)
	if !ok {
		return "", fmt.Errorf("no active %s account for %s at %s", currency, humanID, bankID)
	}

	human := w.Humans[humanID]
	account := w.Accounts[accountID]
	bankCash := bank.CashVault[currency]

	if account.BookedBalance < amount {
		return "", fmt.Errorf(
			"not enough money in %s's account. Available: %d %s",
			humanID,
			account.BookedBalance,
			currency,
		)
	}

	if bankCash < amount {
		return "", fmt.Errorf(
			"not enough physical cash in %s's vault. Available: %d %s",
			bankID,
			bankCash,
			currency,
		)
	}

	account.BookedBalance -= amount
	bank.CashVault[currency] = bankCash - amount
	human.CashWallet[currency] += amount

	return accountID, nil
}

func SupplyCash(w *world.World, centralBankID string, bankID string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	centralBank, ok := w.CentralBanks[centralBankID]
	if !ok {
		return fmt.Errorf("central bank does not exist: %s", centralBankID)
	}

	bank, ok := w.Banks[bankID]
	if !ok {
		return fmt.Errorf("bank does not exist: %s", bankID)
	}

	currency := centralBank.Currency
	bankReservesAtCentralBank, err := reserveBalance(centralBank, bank, centralBankID, bankID)
	if err != nil {
		return err
	}

	if centralBank.CashVault < amount {
		return fmt.Errorf(
			"not enough physical cash in %s vault. Available: %d %s",
			centralBankID,
			centralBank.CashVault,
			currency,
		)
	}

	if bankReservesAtCentralBank < amount {
		return fmt.Errorf(
			"not enough reserves for %s at %s. Available: %d %s",
			bankID,
			centralBankID,
			bankReservesAtCentralBank,
			currency,
		)
	}

	centralBank.CashVault -= amount
	centralBank.ReserveAccounts[bankID] -= amount
	bank.CashVault[currency] += amount
	bank.ReserveBalances[centralBankID] -= amount

	return nil
}

func MoveCash(w *world.World, sourceBankID string, targetBankID string, currency string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if sourceBankID == targetBankID {
		return fmt.Errorf("source bank and target bank must be different")
	}

	sourceBank, ok := w.Banks[sourceBankID]
	if !ok {
		return fmt.Errorf("source bank does not exist: %s", sourceBankID)
	}

	targetBank, ok := w.Banks[targetBankID]
	if !ok {
		return fmt.Errorf("target bank does not exist: %s", targetBankID)
	}

	currency = strings.ToUpper(currency)
	sourceCash := sourceBank.CashVault[currency]

	if sourceCash < amount {
		return fmt.Errorf(
			"not enough physical cash in %s's vault. Available: %d %s",
			sourceBankID,
			sourceCash,
			currency,
		)
	}

	sourceBank.CashVault[currency] = sourceCash - amount
	targetBank.CashVault[currency] += amount

	return nil
}

func SellCash(w *world.World, centralBankID string, sellerBankID string, buyerBankID string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if sellerBankID == buyerBankID {
		return fmt.Errorf("seller bank and buyer bank must be different")
	}

	centralBank, ok := w.CentralBanks[centralBankID]
	if !ok {
		return fmt.Errorf("central bank does not exist: %s", centralBankID)
	}

	sellerBank, ok := w.Banks[sellerBankID]
	if !ok {
		return fmt.Errorf("seller bank does not exist: %s", sellerBankID)
	}

	buyerBank, ok := w.Banks[buyerBankID]
	if !ok {
		return fmt.Errorf("buyer bank does not exist: %s", buyerBankID)
	}

	currency := centralBank.Currency
	sellerCash := sellerBank.CashVault[currency]
	if sellerCash < amount {
		return fmt.Errorf(
			"not enough physical cash in %s's vault. Available: %d %s",
			sellerBankID,
			sellerCash,
			currency,
		)
	}

	if _, err := reserveBalance(centralBank, sellerBank, centralBankID, sellerBankID); err != nil {
		return err
	}

	buyerReserves, err := reserveBalance(centralBank, buyerBank, centralBankID, buyerBankID)
	if err != nil {
		return err
	}

	if buyerReserves < amount {
		return fmt.Errorf(
			"not enough reserves for %s at %s. Available: %d %s",
			buyerBankID,
			centralBankID,
			buyerReserves,
			currency,
		)
	}

	sellerBank.CashVault[currency] = sellerCash - amount
	buyerBank.CashVault[currency] += amount
	centralBank.ReserveAccounts[sellerBankID] += amount
	sellerBank.ReserveBalances[centralBankID] += amount
	centralBank.ReserveAccounts[buyerBankID] -= amount
	buyerBank.ReserveBalances[centralBankID] -= amount

	return nil
}

func ReturnCash(w *world.World, centralBankID string, bankID string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	centralBank, ok := w.CentralBanks[centralBankID]
	if !ok {
		return fmt.Errorf("central bank does not exist: %s", centralBankID)
	}

	bank, ok := w.Banks[bankID]
	if !ok {
		return fmt.Errorf("bank does not exist: %s", bankID)
	}

	currency := centralBank.Currency
	bankCash := bank.CashVault[currency]

	if bankCash < amount {
		return fmt.Errorf(
			"not enough physical cash in %s's vault. Available: %d %s",
			bankID,
			bankCash,
			currency,
		)
	}

	centralBankReserves, err := reserveBalance(centralBank, bank, centralBankID, bankID)
	if err != nil {
		return err
	}

	bank.CashVault[currency] = bankCash - amount
	centralBank.CashVault += amount
	centralBank.ReserveAccounts[bankID] = centralBankReserves + amount
	bank.ReserveBalances[centralBankID] = centralBankReserves + amount

	return nil
}

func DestroyCash(w *world.World, centralBankID string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	centralBank, ok := w.CentralBanks[centralBankID]
	if !ok {
		return fmt.Errorf("central bank does not exist: %s", centralBankID)
	}

	currency := centralBank.Currency

	if centralBank.CashVault < amount {
		return fmt.Errorf(
			"not enough physical cash in %s's vault. Available: %d %s",
			centralBankID,
			centralBank.CashVault,
			currency,
		)
	}

	centralBank.CashVault -= amount
	centralBank.CashIssued -= amount

	return nil
}

func reserveBalance(centralBank *world.CentralBank, bank *world.Bank, centralBankID string, bankID string) (int, error) {
	centralBankReserves, centralBankHasAccount := centralBank.ReserveAccounts[bankID]
	bankReserveMirror, bankHasAccount := bank.ReserveBalances[centralBankID]

	if !centralBankHasAccount || !bankHasAccount {
		return 0, fmt.Errorf("reserve account does not exist: %s at %s", bankID, centralBankID)
	}

	if centralBankReserves != bankReserveMirror {
		return 0, fmt.Errorf(
			"reserve mirror mismatch. %s.reserve_accounts[%s] = %d, %s.reserve_balances[%s] = %d",
			centralBankID,
			bankID,
			centralBankReserves,
			bankID,
			centralBankID,
			bankReserveMirror,
		)
	}

	return centralBankReserves, nil
}
