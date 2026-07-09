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

func newSellCashCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sell-cash central_bank_id seller_bank_id buyer_bank_id amount",
		Short: "Sell physical cash from one commercial bank to another, settled with reserves.",
		Args:  cobra.ExactArgs(4),
		RunE:  runSellCash,
	}
}

func runSellCash(cmd *cobra.Command, args []string) error {
	centralBankID := args[0]
	sellerBankID := args[1]
	buyerBankID := args[2]
	amount, err := strconv.Atoi(args[3])
	if err != nil {
		return fmt.Errorf("amount must be an integer")
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	if err := domain.SellCash(w, centralBankID, sellerBankID, buyerBankID, amount); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "sell-cash", args); err != nil {
		return err
	}

	centralBank := w.CentralBanks[centralBankID]
	sellerBank := w.Banks[sellerBankID]
	buyerBank := w.Banks[buyerBankID]
	currency := centralBank.Currency
	commandlog.Action("Sold %d %s cash from %s to %s", amount, currency, sellerBankID, buyerBankID)
	commandlog.State("%s cash_vault[%s]: %d %s", sellerBankID, currency, sellerBank.CashVault[currency], currency)
	commandlog.State("%s cash_vault[%s]: %d %s", buyerBankID, currency, buyerBank.CashVault[currency], currency)
	commandlog.State("%s reserves at %s: %d %s", sellerBankID, centralBankID, sellerBank.ReserveBalances[centralBankID], currency)
	commandlog.State("%s reserves at %s: %d %s", buyerBankID, centralBankID, buyerBank.ReserveBalances[centralBankID], currency)
	commandlog.State("%s reserve account for %s: %d %s", centralBankID, sellerBankID, centralBank.ReserveAccounts[sellerBankID], currency)
	commandlog.State("%s reserve account for %s: %d %s", centralBankID, buyerBankID, centralBank.ReserveAccounts[buyerBankID], currency)
	return nil
}
