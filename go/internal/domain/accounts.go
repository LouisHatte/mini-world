package domain

import "mini-world-go/internal/world"

func ActiveAccountID(w *world.World, humanID, bankID, currency string) (string, bool) {
	accountID := BuildAccountID(bankID, humanID, currency)
	account, ok := w.Accounts[accountID]

	if !ok || account.Status != world.AccountActive {
		return "", false
	}

	return accountID, true
}
