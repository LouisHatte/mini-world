package domain

import (
	"fmt"
	"strings"

	"mini-world-go/internal/world"
)

func CreateCentralBank(w *world.World, centralBankID string, currency string) error {
	if _, exists := w.CentralBanks[centralBankID]; exists {
		return fmt.Errorf("central bank already exists: %s", centralBankID)
	}

	currency = strings.ToUpper(currency)
	w.CentralBanks[centralBankID] = world.NewCentralBank(centralBankID, currency)

	return nil
}

func CreateBank(w *world.World, bankID string) error {
	if _, exists := w.Banks[bankID]; exists {
		return fmt.Errorf("bank already exists: %s", bankID)
	}

	w.Banks[bankID] = world.NewBank(bankID)

	return nil
}

func CreateHuman(w *world.World, humanID string) error {
	if _, exists := w.Humans[humanID]; exists {
		return fmt.Errorf("human already exists: %s", humanID)
	}

	w.Humans[humanID] = world.NewHuman(humanID)

	return nil
}

func OpenAccount(w *world.World, humanID string, bankID string, currency string) (string, error) {
	human, humanExists := w.Humans[humanID]
	if !humanExists {
		return "", fmt.Errorf("human does not exist: %s", humanID)
	}

	bank, bankExists := w.Banks[bankID]
	if !bankExists {
		return "", fmt.Errorf("bank does not exist: %s", bankID)
	}

	currency = strings.ToUpper(currency)
	accountID := BuildAccountID(bankID, humanID, currency)

	if _, exists := w.Accounts[accountID]; exists {
		return "", fmt.Errorf("account already exists: %s", accountID)
	}

	w.Accounts[accountID] = world.NewAccount(accountID, humanID, bankID, currency)

	human.BankAccounts = append(human.BankAccounts, accountID)
	bank.CustomerAccounts = append(bank.CustomerAccounts, accountID)

	return accountID, nil
}

func OpenReserveAccount(w *world.World, centralBankID string, bankID string) error {
	centralBank, centralBankExists := w.CentralBanks[centralBankID]
	if !centralBankExists {
		return fmt.Errorf("central bank does not exist: %s", centralBankID)
	}

	bank, bankExists := w.Banks[bankID]
	if !bankExists {
		return fmt.Errorf("bank does not exist: %s", bankID)
	}

	if _, exists := centralBank.ReserveAccounts[bankID]; exists {
		return fmt.Errorf("reserve account already exists: %s at %s", bankID, centralBankID)
	}

	if _, exists := bank.ReserveBalances[centralBankID]; exists {
		return fmt.Errorf("reserve balance already exists: %s at %s", bankID, centralBankID)
	}

	centralBank.ReserveAccounts[bankID] = 0
	bank.ReserveBalances[centralBankID] = 0

	return nil
}

func BuildAccountID(bankID string, humanID string, currency string) string {
	return fmt.Sprintf("acc_%s_%s_%s", bankID, humanID, strings.ToLower(currency))
}
