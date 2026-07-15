package swift

import (
	"strings"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newOpenCorrespondentAccountCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "open-correspondent-account owner_bank_id correspondent_bank_id currency",
		Short: "Open a correspondent account between two banks.",
		Args:  cobra.ExactArgs(3),
		RunE:  runOpenCorrespondentAccount,
	}
}

func runOpenCorrespondentAccount(cmd *cobra.Command, args []string) error {
	ownerBankID := args[0]
	correspondentBankID := args[1]
	currency := strings.ToUpper(args[2])

	w, err := world.Load()
	if err != nil {
		return err
	}

	accountID, err := domain.OpenCorrespondentAccount(w, ownerBankID, correspondentBankID, currency)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "open-correspondent-account", args); err != nil {
		return err
	}

	commandlog.Action("Opened correspondent account: %s", accountID)
	commandlog.State("Owner bank: %s", ownerBankID)
	commandlog.State("Correspondent bank: %s", correspondentBankID)
	commandlog.State("Currency: %s", currency)
	logCorrespondentAccount(w, accountID)
	return nil
}
