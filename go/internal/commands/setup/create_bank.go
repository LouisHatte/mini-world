package setup

import (
	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newCreateBankCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create-bank bank_id",
		Short: "Create a commercial bank.",
		Args:  cobra.ExactArgs(1),
		RunE:  runCreateBank,
	}
}

func runCreateBank(cmd *cobra.Command, args []string) error {
	bankID := args[0]

	w, err := world.Load()
	if err != nil {
		return err
	}

	if err := domain.CreateBank(w, bankID); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "create-bank", args); err != nil {
		return err
	}

	commandlog.Action("Created bank: %s", bankID)
	return nil
}
