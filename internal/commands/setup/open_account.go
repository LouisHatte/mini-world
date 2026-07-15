package setup

import (
	"strings"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newOpenAccountCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "open-account human_id bank_id currency",
		Short: "Open a bank account for a human at a commercial bank.",
		Args:  cobra.ExactArgs(3),
		RunE:  runOpenAccount,
	}
}

func runOpenAccount(cmd *cobra.Command, args []string) error {
	humanID := args[0]
	bankID := args[1]
	currency := strings.ToUpper(args[2])

	w, err := world.Load()
	if err != nil {
		return err
	}

	accountID, err := domain.OpenAccount(w, humanID, bankID, currency)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "open-account", args); err != nil {
		return err
	}

	commandlog.Action("Opened account: %s", accountID)
	commandlog.State("Owner: %s", humanID)
	commandlog.State("Bank: %s", bankID)
	commandlog.State("Currency: %s", currency)
	commandlog.State("Booked balance: 0")
	return nil
}
