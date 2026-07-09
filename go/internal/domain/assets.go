package domain

import (
	"fmt"
	"strings"

	"mini-world-go/internal/world"
)

func RegisterAsset(w *world.World, assetID string, ownerHumanID string, currency string, estimatedValue int) error {
	if estimatedValue < 0 {
		return fmt.Errorf("estimated value must be greater than or equal to 0")
	}

	if _, exists := w.Assets[assetID]; exists {
		return fmt.Errorf("asset already exists: %s", assetID)
	}

	if _, exists := w.Humans[ownerHumanID]; !exists {
		return fmt.Errorf("owner human does not exist: %s", ownerHumanID)
	}

	currency = strings.ToUpper(currency)
	if _, exists := w.Currencies[currency]; !exists {
		return fmt.Errorf("currency does not exist: %s", currency)
	}

	w.Assets[assetID] = world.NewAsset(assetID, ownerHumanID, currency, estimatedValue)

	return nil
}

func BuyAssetCash(w *world.World, buyerHumanID string, assetID string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	buyer, ok := w.Humans[buyerHumanID]
	if !ok {
		return fmt.Errorf("buyer human does not exist: %s", buyerHumanID)
	}

	asset, ok := w.Assets[assetID]
	if !ok {
		return fmt.Errorf("asset does not exist: %s", assetID)
	}

	if asset.OwnerHumanID == buyerHumanID {
		return fmt.Errorf("%s already owns asset: %s", buyerHumanID, assetID)
	}

	seller, ok := w.Humans[asset.OwnerHumanID]
	if !ok {
		return fmt.Errorf("asset owner human does not exist: %s", asset.OwnerHumanID)
	}

	currency := asset.Currency
	buyerCash := buyer.CashWallet[currency]
	if buyerCash < amount {
		return fmt.Errorf(
			"not enough cash in %s's wallet. Available: %d %s",
			buyerHumanID,
			buyerCash,
			currency,
		)
	}

	buyer.CashWallet[currency] = buyerCash - amount
	seller.CashWallet[currency] += amount
	asset.OwnerHumanID = buyerHumanID

	return nil
}

func RevalueAsset(w *world.World, assetID string, estimatedValue int) error {
	if estimatedValue < 0 {
		return fmt.Errorf("estimated value must be greater than or equal to 0")
	}

	asset, ok := w.Assets[assetID]
	if !ok {
		return fmt.Errorf("asset does not exist: %s", assetID)
	}

	asset.EstimatedValue = estimatedValue

	return nil
}
