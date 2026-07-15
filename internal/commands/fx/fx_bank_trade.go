package fx

import (
	"strings"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newFXBankTradeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "fx-bank-trade from_bank_id to_bank_id from_currency to_currency amount",
		Short: "Two banks exchange currencies, settled with reserves.",
		Args:  cobra.ExactArgs(5),
		RunE:  runFXBankTrade,
	}
}

func runFXBankTrade(cmd *cobra.Command, args []string) error {
	fromBankID := args[0]
	toBankID := args[1]
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

	result, err := domain.FXBankTrade(w, fromBankID, toBankID, fromCurrency, toCurrency, amount)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "fx-bank-trade", args); err != nil {
		return err
	}

	fromCentralBank := w.CentralBanks[result.FromCentralBankID]
	toCentralBank := w.CentralBanks[result.ToCentralBankID]
	fromBank := w.Banks[fromBankID]
	toBank := w.Banks[toBankID]
	commandlog.Action("Traded FX between %s and %s: %d %s -> %d %s", fromBankID, toBankID, result.FromAmount, fromCurrency, result.ToAmount, toCurrency)
	commandlog.State("Rate: %.6f", result.Rate)
	commandlog.State("%s reserve account for %s: %d %s", result.FromCentralBankID, fromBankID, fromCentralBank.ReserveAccounts[fromBankID], fromCurrency)
	commandlog.State("%s reserve account for %s: %d %s", result.FromCentralBankID, toBankID, fromCentralBank.ReserveAccounts[toBankID], fromCurrency)
	commandlog.State("%s reserve account for %s: %d %s", result.ToCentralBankID, fromBankID, toCentralBank.ReserveAccounts[fromBankID], toCurrency)
	commandlog.State("%s reserve account for %s: %d %s", result.ToCentralBankID, toBankID, toCentralBank.ReserveAccounts[toBankID], toCurrency)
	commandlog.State("%s reserves at %s: %d %s", fromBankID, result.FromCentralBankID, fromBank.ReserveBalances[result.FromCentralBankID], fromCurrency)
	commandlog.State("%s reserves at %s: %d %s", toBankID, result.FromCentralBankID, toBank.ReserveBalances[result.FromCentralBankID], fromCurrency)
	commandlog.State("%s reserves at %s: %d %s", fromBankID, result.ToCentralBankID, fromBank.ReserveBalances[result.ToCentralBankID], toCurrency)
	commandlog.State("%s reserves at %s: %d %s", toBankID, result.ToCentralBankID, toBank.ReserveBalances[result.ToCentralBankID], toCurrency)
	return nil
}
