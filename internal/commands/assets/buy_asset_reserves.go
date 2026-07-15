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

func newBuyAssetReservesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "buy-asset-reserves central_bank_id buyer_id asset_id amount",
		Short: "Buy an asset between banks or central banks, settled with reserves.",
		Args:  cobra.ExactArgs(4),
		RunE:  runBuyAssetReserves,
	}
}

func runBuyAssetReserves(cmd *cobra.Command, args []string) error {
	centralBankID := args[0]
	buyerID := args[1]
	assetID := args[2]
	amount, err := strconv.Atoi(args[3])
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
	sellerType := asset.OwnerType
	sellerID := asset.OwnerID
	currency := asset.Currency

	if err := domain.BuyAssetReserves(w, centralBankID, buyerID, assetID, amount); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "buy-asset-reserves", args); err != nil {
		return err
	}

	asset = w.Assets[assetID]
	commandlog.Action("%s bought asset %s from %s for %d %s reserves", buyerID, assetID, sellerID, amount, currency)
	commandlog.State("%s owner: %s %s", assetID, asset.OwnerType, asset.OwnerID)
	logReserveSettlement(w, centralBankID, sellerType, sellerID, asset.OwnerType, buyerID, currency)
	return nil
}

func logReserveSettlement(w *world.World, centralBankID string, sellerType world.AssetOwnerType, sellerID string, buyerType world.AssetOwnerType, buyerID string, currency string) {
	centralBank := w.CentralBanks[centralBankID]

	if sellerType == world.AssetOwnerBank {
		sellerBank := w.Banks[sellerID]
		commandlog.State("%s reserve account for %s: %d %s", centralBankID, sellerID, centralBank.ReserveAccounts[sellerID], currency)
		commandlog.State("%s reserves at %s: %d %s", sellerID, centralBankID, sellerBank.ReserveBalances[centralBankID], currency)
	}

	if buyerType == world.AssetOwnerBank {
		buyerBank := w.Banks[buyerID]
		commandlog.State("%s reserve account for %s: %d %s", centralBankID, buyerID, centralBank.ReserveAccounts[buyerID], currency)
		commandlog.State("%s reserves at %s: %d %s", buyerID, centralBankID, buyerBank.ReserveBalances[centralBankID], currency)
	}
}
