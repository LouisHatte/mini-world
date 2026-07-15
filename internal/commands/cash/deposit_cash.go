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

func newDepositCashCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "deposit-cash human_id bank_id currency amount",
		Short: "Deposit physical cash from a human wallet into a bank account.",
		Args:  cobra.ExactArgs(4),
		RunE:  runDepositCash,
	}
}

func runDepositCash(cmd *cobra.Command, args []string) error {
	humanID := args[0]
	bankID := args[1]
	currency := strings.ToUpper(args[2])
	amount, err := strconv.Atoi(args[3])
	if err != nil {
		return fmt.Errorf("amount must be an integer")
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	accountID, err := domain.DepositCash(w, humanID, bankID, currency, amount)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "deposit-cash", args); err != nil {
		return err
	}

	human := w.Humans[humanID]
	bank := w.Banks[bankID]
	account := w.Accounts[accountID]
	commandlog.Action("%s deposited %d %s cash at %s", humanID, amount, currency, bankID)
	commandlog.State("%s cash_wallet[%s]: %d %s", humanID, currency, human.CashWallet[currency], currency)
	commandlog.State("%s cash_vault[%s]: %d %s", bankID, currency, bank.CashVault[currency], currency)
	commandlog.State("%s booked_balance: %d %s", accountID, account.BookedBalance, currency)
	return nil
}
