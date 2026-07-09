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

func newTransferCashCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "transfer-cash source_human_id target_human_id currency amount",
		Short: "Transfer physical cash from one human to another.",
		Args:  cobra.ExactArgs(4),
		RunE:  runTransferCash,
	}
}

func runTransferCash(cmd *cobra.Command, args []string) error {
	sourceHumanID := args[0]
	targetHumanID := args[1]
	currency := strings.ToUpper(args[2])
	amount, err := strconv.Atoi(args[3])
	if err != nil {
		return fmt.Errorf("amount must be an integer")
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	if err := domain.TransferCash(w, sourceHumanID, targetHumanID, currency, amount); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "transfer-cash", args); err != nil {
		return err
	}

	sourceHuman := w.Humans[sourceHumanID]
	targetHuman := w.Humans[targetHumanID]
	commandlog.Action("Transferred %d %s cash from %s to %s", amount, currency, sourceHumanID, targetHumanID)
	commandlog.State("%s cash_wallet[%s]: %d %s", sourceHumanID, currency, sourceHuman.CashWallet[currency], currency)
	commandlog.State("%s cash_wallet[%s]: %d %s", targetHumanID, currency, targetHuman.CashWallet[currency], currency)
	return nil
}
