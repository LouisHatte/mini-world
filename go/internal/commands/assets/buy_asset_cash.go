package assets

import (
	"fmt"
	"strconv"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newBuyAssetCashCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "buy-asset-cash buyer_human_id asset_id amount",
		Short: "Buy an existing asset with physical cash.",
		Args:  cobra.ExactArgs(3),
		RunE:  runBuyAssetCash,
	}
}

func runBuyAssetCash(cmd *cobra.Command, args []string) error {
	buyerHumanID := args[0]
	assetID := args[1]
	amount, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("amount must be an integer")
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	asset, ok := w.Assets[assetID]
	if !ok {
		return commandrun.PrintBusinessError(fmt.Errorf("asset does not exist: %s", assetID))
	}
	sellerHumanID := asset.OwnerHumanID
	currency := asset.Currency

	if err := domain.BuyAssetCash(w, buyerHumanID, assetID, amount); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "buy-asset-cash", args); err != nil {
		return err
	}

	buyer := w.Humans[buyerHumanID]
	seller := w.Humans[sellerHumanID]
	commandlog.Action("%s bought asset %s from %s for %d %s cash", buyerHumanID, assetID, sellerHumanID, amount, currency)
	commandlog.State("%s owner: %s", assetID, w.Assets[assetID].OwnerHumanID)
	commandlog.State("%s cash_wallet[%s]: %d %s", buyerHumanID, currency, buyer.CashWallet[currency], currency)
	commandlog.State("%s cash_wallet[%s]: %d %s", sellerHumanID, currency, seller.CashWallet[currency], currency)
	return nil
}
