package reserves

import (
	"fmt"
	"strconv"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newRepayReserveLoanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "repay-reserve-loan bank_id reserve_loan_id amount",
		Short: "Commercial bank repays a central bank reserve loan.",
		Args:  cobra.ExactArgs(3),
		RunE:  runRepayReserveLoan,
	}
}

func runRepayReserveLoan(cmd *cobra.Command, args []string) error {
	bankID := args[0]
	reserveLoanID := args[1]
	amount, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("amount must be an integer")
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	reserveLoan, ok := w.ReserveLoans[reserveLoanID]
	if !ok {
		return commandrun.PrintBusinessError(fmt.Errorf("reserve loan does not exist: %s", reserveLoanID))
	}
	centralBankID := reserveLoan.CentralBankID
	currency := reserveLoan.Currency

	if err := domain.RepayReserveLoan(w, bankID, reserveLoanID, amount); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "repay-reserve-loan", args); err != nil {
		return err
	}

	centralBank := w.CentralBanks[centralBankID]
	bank := w.Banks[bankID]
	commandlog.Action("%s repaid %d %s on reserve loan %s", bankID, amount, currency, reserveLoanID)
	commandlog.State("%s reserve account for %s: %d %s", centralBankID, bankID, centralBank.ReserveAccounts[bankID], currency)
	commandlog.State("%s reserves at %s: %d %s", bankID, centralBankID, bank.ReserveBalances[centralBankID], currency)
	commandlog.State("%s loan to %s: %d %s", centralBankID, bankID, centralBank.LoansToBanks[bankID], currency)
	commandlog.State("%s loan from %s: %d %s", bankID, centralBankID, bank.LoansFromCentralBanks[centralBankID], currency)
	commandlog.State("%s outstanding: %d %s", reserveLoanID, w.ReserveLoans[reserveLoanID].Outstanding, currency)
	commandlog.State("%s status: %s", reserveLoanID, w.ReserveLoans[reserveLoanID].Status)
	return nil
}
