package loans

import (
	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newDefaultLoanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "default-loan loan_id",
		Short: "Mark a customer loan as defaulted.",
		Args:  cobra.ExactArgs(1),
		RunE:  runDefaultLoan,
	}
}

func runDefaultLoan(cmd *cobra.Command, args []string) error {
	loanID := args[0]

	w, err := world.Load()
	if err != nil {
		return err
	}

	if err := domain.DefaultLoan(w, loanID); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "default-loan", args); err != nil {
		return err
	}

	loan := w.CustomerLoans[loanID]
	totalDue := loan.OutstandingPrincipal + loan.OutstandingInterest
	commandlog.Action("Defaulted loan: %s", loanID)
	logLoanState(w, loanID)
	commandlog.State("%s total_due: %d %s", loanID, totalDue, loan.Currency)
	return nil
}
