package sepa

import (
	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newSettleSepaCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "settle-sepa payment_id",
		Short: "Settle a SEPA payment with central bank reserves.",
		Args:  cobra.ExactArgs(1),
		RunE:  runSettleSepa,
	}
}

func runSettleSepa(cmd *cobra.Command, args []string) error {
	paymentID := args[0]

	w, err := world.Load()
	if err != nil {
		return err
	}

	result, err := domain.SettleSepa(w, paymentID)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "settle-sepa", args); err != nil {
		return err
	}

	if result.SettlementFailed {
		commandlog.Action("SEPA settlement failed: %s", paymentID)
		logSepaSettlementResult(w, result)
		return nil
	}

	commandlog.Action("Settled SEPA payment: %s", paymentID)
	logSepaSettlementResult(w, result)
	return nil
}
