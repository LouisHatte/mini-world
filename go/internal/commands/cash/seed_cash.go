package cash

import (
	"fmt"
	"strconv"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newSeedCashCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "seed-cash central_bank_id human_id amount",
		Short: "Seed already-issued physical cash from a central bank vault to a human.",
		Args:  cobra.ExactArgs(3),
		RunE:  runSeedCash,
	}
}

func runSeedCash(cmd *cobra.Command, args []string) error {
	centralBankID := args[0]
	humanID := args[1]
	amount, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("amount must be an integer")
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	if err := domain.SeedCash(w, centralBankID, humanID, amount); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "seed-cash", args); err != nil {
		return err
	}

	centralBank := w.CentralBanks[centralBankID]
	human := w.Humans[humanID]
	currency := centralBank.Currency
	commandlog.Action("Seeded %d %s cash from %s to %s", amount, currency, centralBankID, humanID)
	commandlog.State("%s cash_vault: %d %s", centralBankID, centralBank.CashVault, currency)
	commandlog.State("%s cash_wallet[%s]: %d %s", humanID, currency, human.CashWallet[currency], currency)
	return nil
}
