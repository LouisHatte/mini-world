package domain

import (
	"fmt"
	"strings"

	"mini-world-go/internal/world"
)

func RegisterAsset(w *world.World, assetID string, ownerID string, currency string, estimatedValue int) error {
	if estimatedValue < 0 {
		return fmt.Errorf("estimated value must be greater than or equal to 0")
	}

	if _, exists := w.Assets[assetID]; exists {
		return fmt.Errorf("asset already exists: %s", assetID)
	}

	ownerType, err := assetOwnerType(w, ownerID)
	if err != nil {
		return err
	}

	currency = strings.ToUpper(currency)
	if _, exists := w.Currencies[currency]; !exists {
		return fmt.Errorf("currency does not exist: %s", currency)
	}

	w.Assets[assetID] = world.NewAsset(assetID, ownerType, ownerID, currency, estimatedValue)

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

	if asset.OwnerType != world.AssetOwnerHuman {
		return fmt.Errorf("asset owner is not a human: %s", asset.OwnerID)
	}

	if asset.OwnerID == buyerHumanID {
		return fmt.Errorf("%s already owns asset: %s", buyerHumanID, assetID)
	}

	seller, ok := w.Humans[asset.OwnerID]
	if !ok {
		return fmt.Errorf("asset owner human does not exist: %s", asset.OwnerID)
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
	asset.OwnerType = world.AssetOwnerHuman
	asset.OwnerID = buyerHumanID

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

func BuyAssetReserves(w *world.World, centralBankID string, buyerID string, assetID string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	centralBank, ok := w.CentralBanks[centralBankID]
	if !ok {
		return fmt.Errorf("central bank does not exist: %s", centralBankID)
	}

	asset, ok := w.Assets[assetID]
	if !ok {
		return fmt.Errorf("asset does not exist: %s", assetID)
	}

	if asset.Currency != centralBank.Currency {
		return fmt.Errorf("asset currency is %s, not %s", asset.Currency, centralBank.Currency)
	}

	if asset.CollateralForReserveLoanID != "" {
		return fmt.Errorf("asset is pledged to reserve loan: %s", asset.CollateralForReserveLoanID)
	}

	sellerID := asset.OwnerID
	sellerType, err := assetOwnerType(w, sellerID)
	if err != nil {
		return err
	}

	buyerType, err := assetOwnerType(w, buyerID)
	if err != nil {
		return err
	}

	if sellerType == world.AssetOwnerHuman || buyerType == world.AssetOwnerHuman {
		return fmt.Errorf("humans cannot settle assets with reserves")
	}

	if sellerID == buyerID {
		return fmt.Errorf("%s already owns asset: %s", buyerID, assetID)
	}

	if sellerType == world.AssetOwnerCentralBank && buyerType == world.AssetOwnerCentralBank {
		return fmt.Errorf("central bank to central bank asset transfer is not supported")
	}

	if asset.OwnerType != sellerType || asset.OwnerID != sellerID {
		return fmt.Errorf("asset %s is owned by %s %s, not %s %s", assetID, asset.OwnerType, asset.OwnerID, sellerType, sellerID)
	}

	if buyerType == world.AssetOwnerBank {
		buyerBank := w.Banks[buyerID]
		buyerReserves, err := reserveBalance(centralBank, buyerBank, centralBankID, buyerID)
		if err != nil {
			return err
		}

		if buyerReserves < amount {
			return fmt.Errorf(
				"not enough reserves for %s at %s. Available: %d %s",
				buyerID,
				centralBankID,
				buyerReserves,
				centralBank.Currency,
			)
		}

		centralBank.ReserveAccounts[buyerID] -= amount
		buyerBank.ReserveBalances[centralBankID] -= amount
	}

	if sellerType == world.AssetOwnerBank {
		sellerBank := w.Banks[sellerID]
		if _, err := reserveBalance(centralBank, sellerBank, centralBankID, sellerID); err != nil {
			return err
		}

		centralBank.ReserveAccounts[sellerID] += amount
		sellerBank.ReserveBalances[centralBankID] += amount
	}

	asset.OwnerType = buyerType
	asset.OwnerID = buyerID

	return nil
}

func assetOwnerType(w *world.World, ownerID string) (world.AssetOwnerType, error) {
	_, humanExists := w.Humans[ownerID]
	_, bankExists := w.Banks[ownerID]
	_, centralBankExists := w.CentralBanks[ownerID]

	if matches := boolAsInt(humanExists) + boolAsInt(bankExists) + boolAsInt(centralBankExists); matches > 1 {
		return "", fmt.Errorf("asset owner id is ambiguous: %s", ownerID)
	}

	if humanExists {
		return world.AssetOwnerHuman, nil
	}

	if bankExists {
		return world.AssetOwnerBank, nil
	}

	if centralBankExists {
		return world.AssetOwnerCentralBank, nil
	}

	return "", fmt.Errorf("asset owner does not exist: %s", ownerID)
}

func boolAsInt(value bool) int {
	if value {
		return 1
	}

	return 0
}
