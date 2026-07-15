package fx

import (
	"strings"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newSetFXRateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set-fx-rate from_currency to_currency rate",
		Short: "Set exchange rate between two currencies.",
		Args:  cobra.ExactArgs(3),
		RunE:  runSetFXRate,
	}
}

func runSetFXRate(cmd *cobra.Command, args []string) error {
	fromCurrency := strings.ToUpper(args[0])
	toCurrency := strings.ToUpper(args[1])
	rate, err := parseRate(args[2])
	if err != nil {
		return err
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	marketID, err := domain.SetFXRate(w, fromCurrency, toCurrency, rate)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "set-fx-rate", args); err != nil {
		return err
	}

	commandlog.Action("Set FX rate: %s", marketID)
	commandlog.State("From: %s", fromCurrency)
	commandlog.State("To: %s", toCurrency)
	commandlog.State("Rate: %.6f", rate)
	return nil
}
