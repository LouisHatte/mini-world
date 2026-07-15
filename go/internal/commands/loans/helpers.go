package loans

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

func logLoanState(w *world.World, loanID string) {
	loan := w.CustomerLoans[loanID]

	commandlog.State("%s outstanding_principal: %d %s", loanID, loan.OutstandingPrincipal, loan.Currency)
	commandlog.State("%s outstanding_interest: %d %s", loanID, loan.OutstandingInterest, loan.Currency)
	commandlog.State("%s status: %s", loanID, loan.Status)
	commandlog.State("%s loans[%s]: %d %s", loan.BorrowerHumanID, loanID, w.Humans[loan.BorrowerHumanID].Loans[loanID], loan.Currency)
}
