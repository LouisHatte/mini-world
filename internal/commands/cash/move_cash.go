package cash

import (
	"fmt"
	"strconv"
	"strings"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newMoveCashCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "move-cash source_bank_id target_bank_id currency amount",
		Short: "Move physical cash from one commercial bank vault to another.",
		Args:  cobra.ExactArgs(4),
		RunE:  runMoveCash,
	}
}

func runMoveCash(cmd *cobra.Command, args []string) error {
	sourceBankID := args[0]
	targetBankID := args[1]
	currency := strings.ToUpper(args[2])
	amount, err := strconv.Atoi(args[3])
	if err != nil {
		return fmt.Errorf("amount must be an integer")
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	if err := domain.MoveCash(w, sourceBankID, targetBankID, currency, amount); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "move-cash", args); err != nil {
		return err
	}

	sourceBank := w.Banks[sourceBankID]
	targetBank := w.Banks[targetBankID]
	commandlog.Action("Moved %d %s cash from %s to %s", amount, currency, sourceBankID, targetBankID)
	commandlog.State("%s cash_vault[%s]: %d %s", sourceBankID, currency, sourceBank.CashVault[currency], currency)
	commandlog.State("%s cash_vault[%s]: %d %s", targetBankID, currency, targetBank.CashVault[currency], currency)
	return nil
}
