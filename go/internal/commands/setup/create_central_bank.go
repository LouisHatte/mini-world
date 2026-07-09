package setup

import (
	"strings"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newCreateCentralBankCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create-central-bank central_bank_id currency",
		Short: "Create a central bank.",
		Args:  cobra.ExactArgs(2),
		RunE:  runCreateCentralBank,
	}
}

func runCreateCentralBank(cmd *cobra.Command, args []string) error {
	centralBankID := args[0]
	currency := strings.ToUpper(args[1])

	w, err := world.Load()
	if err != nil {
		return err
	}

	if err := domain.CreateCentralBank(w, centralBankID, currency); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "create-central-bank", args); err != nil {
		return err
	}

	commandlog.Action("Created central bank: %s (%s)", centralBankID, currency)
	return nil
}
