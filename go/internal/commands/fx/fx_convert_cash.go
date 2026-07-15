package fx

import (
	"strings"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newFXConvertCashCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "fx-convert-cash human_id bank_id from_currency to_currency amount",
		Short: "Human exchanges physical cash from one currency to another through a bank.",
		Args:  cobra.ExactArgs(5),
		RunE:  runFXConvertCash,
	}
}

func runFXConvertCash(cmd *cobra.Command, args []string) error {
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

	result, err := domain.FXConvertCash(w, humanID, bankID, fromCurrency, toCurrency, amount)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "fx-convert-cash", args); err != nil {
		return err
	}

	human := w.Humans[humanID]
	commandlog.Action("Converted cash for %s through %s: %d %s -> %d %s", humanID, bankID, result.FromAmount, fromCurrency, result.ToAmount, toCurrency)
	commandlog.State("Rate: %.6f", result.Rate)
	commandlog.State("%s cash_wallet[%s]: %d %s", humanID, fromCurrency, human.CashWallet[fromCurrency], fromCurrency)
	commandlog.State("%s cash_wallet[%s]: %d %s", humanID, toCurrency, human.CashWallet[toCurrency], toCurrency)
	logBankFXState(w, bankID, fromCurrency, toCurrency)
	return nil
}
