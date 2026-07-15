package fx

import (
	"fmt"
	"strconv"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/world"
)

func parseAmount(value string) (int, error) {
	amount, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("amount must be an integer")
	}

	return amount, nil
}

func parseRate(value string) (float64, error) {
	rate, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("rate must be a number")
	}

	return rate, nil
}

func logBankFXState(w *world.World, bankID string, currencies ...string) {
	bank := w.Banks[bankID]

	for _, currency := range currencies {
		commandlog.State("%s fx_inventory[%s]: %d %s", bankID, currency, bank.FXInventory[currency], currency)
		commandlog.State("%s cash_vault[%s]: %d %s", bankID, currency, bank.CashVault[currency], currency)
	}
}
