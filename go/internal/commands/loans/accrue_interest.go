package loans

import (
	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newAccrueInterestCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "accrue-interest loan_id amount",
		Short: "Add interest to an active customer loan.",
		Args:  cobra.ExactArgs(2),
		RunE:  runAccrueInterest,
	}
}

func runAccrueInterest(cmd *cobra.Command, args []string) error {
	loanID := args[0]
	amount, err := parseAmount(args[1])
	if err != nil {
		return err
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	currency := ""
	if loan, ok := w.CustomerLoans[loanID]; ok {
		currency = loan.Currency
	}

	if err := domain.AccrueInterest(w, loanID, amount); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "accrue-interest", args); err != nil {
		return err
	}

	commandlog.Action("Accrued %d %s interest on %s", amount, currency, loanID)
	logLoanState(w, loanID)
	return nil
}
