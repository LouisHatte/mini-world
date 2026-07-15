package swift

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

func logCorrespondentAccount(w *world.World, correspondentAccountID string) {
	account := w.CorrespondentAccounts[correspondentAccountID]
	commandlog.State("%s nostro_balance: %d %s", correspondentAccountID, account.NostroBalance, account.Currency)
	commandlog.State("%s vostro_balance: %d %s", correspondentAccountID, account.VostroBalance, account.Currency)
}
