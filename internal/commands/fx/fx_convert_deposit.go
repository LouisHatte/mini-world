package fx

import (
	"strings"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newFXConvertDepositCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "fx-convert-deposit human_id bank_id from_currency to_currency amount",
		Short: "Human converts bank deposit from one currency to another through a bank.",
		Args:  cobra.ExactArgs(5),
		RunE:  runFXConvertDeposit,
	}
}

func runFXConvertDeposit(cmd *cobra.Command, args []string) error {
	humanID := args[0]
	bankID := args[1]
	fromCurrency := strings.ToUpper(args[2])
	toCurrency := strings.ToUpper(args[3])
	amount, err := parseAmount(args[4])
	if err != nil {
		return err
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	result, err := domain.FXConvertDeposit(w, humanID, bankID, fromCurrency, toCurrency, amount)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "fx-convert-deposit", args); err != nil {
		return err
	}

	commandlog.Action("Converted deposit for %s at %s: %d %s -> %d %s", humanID, bankID, result.FromAmount, fromCurrency, result.ToAmount, toCurrency)
	commandlog.State("Rate: %.6f", result.Rate)
	commandlog.State("%s booked_balance: %d %s", result.SourceAccountID, w.Accounts[result.SourceAccountID].BookedBalance, fromCurrency)
	commandlog.State("%s booked_balance: %d %s", result.TargetAccountID, w.Accounts[result.TargetAccountID].BookedBalance, toCurrency)
	logBankFXState(w, bankID, fromCurrency, toCurrency)
	return nil
}
