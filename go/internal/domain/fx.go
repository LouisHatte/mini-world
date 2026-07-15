package domain

import (
	"fmt"
	"math"
	"strings"

	"mini-world-go/internal/world"
)

type FXConversionResult struct {
	FromAmount int
	ToAmount   int
	Rate       float64
}

type FXDepositConversionResult struct {
	FXConversionResult
	SourceAccountID string
	TargetAccountID string
}

type FXCashConversionResult struct {
	FXConversionResult
}

type FXBankTradeResult struct {
	FXConversionResult
	FromCentralBankID string
	ToCentralBankID   string
}

func SetFXRate(w *world.World, fromCurrency string, toCurrency string, rate float64) (string, error) {
	fromCurrency = strings.ToUpper(fromCurrency)
	toCurrency = strings.ToUpper(toCurrency)

	if fromCurrency == toCurrency {
		return "", fmt.Errorf("from currency and to currency must be different")
	}

	if rate <= 0 {
		return "", fmt.Errorf("rate must be greater than 0")
	}

	if _, exists := w.Currencies[fromCurrency]; !exists {
		return "", fmt.Errorf("currency does not exist: %s", fromCurrency)
	}

	if _, exists := w.Currencies[toCurrency]; !exists {
		return "", fmt.Errorf("currency does not exist: %s", toCurrency)
	}

	marketID := FXMarketID(fromCurrency, toCurrency)
	w.FXMarkets[marketID] = world.NewFXMarket(marketID, fromCurrency, toCurrency, rate)

	return marketID, nil
}

func FXConvertDeposit(w *world.World, humanID string, bankID string, fromCurrency string, toCurrency string, amount int) (FXDepositConversionResult, error) {
	fromCurrency = strings.ToUpper(fromCurrency)
	toCurrency = strings.ToUpper(toCurrency)

	if amount <= 0 {
		return FXDepositConversionResult{}, fmt.Errorf("amount must be greater than 0")
	}

	bank, ok := w.Banks[bankID]
	if !ok {
		return FXDepositConversionResult{}, fmt.Errorf("bank does not exist: %s", bankID)
	}

	if _, ok := w.Humans[humanID]; !ok {
		return FXDepositConversionResult{}, fmt.Errorf("human does not exist: %s", humanID)
	}

	rate, convertedAmount, err := convertAmount(w, fromCurrency, toCurrency, amount)
	if err != nil {
		return FXDepositConversionResult{}, err
	}

	sourceAccountID, sourceAccount, err := activeAccount(w, humanID, bankID, fromCurrency)
	if err != nil {
		return FXDepositConversionResult{}, err
	}

	targetAccountID, targetAccount, err := activeAccount(w, humanID, bankID, toCurrency)
	if err != nil {
		return FXDepositConversionResult{}, err
	}

	if sourceAccount.BookedBalance < amount {
		return FXDepositConversionResult{}, fmt.Errorf(
			"not enough deposits in %s. Available: %d %s",
			sourceAccountID,
			sourceAccount.BookedBalance,
			fromCurrency,
		)
	}

	sourceAccount.BookedBalance -= amount
	targetAccount.BookedBalance += convertedAmount
	bank.FXInventory[fromCurrency] += amount
	bank.FXInventory[toCurrency] -= convertedAmount

	return FXDepositConversionResult{
		FXConversionResult: FXConversionResult{
			FromAmount: amount,
			ToAmount:   convertedAmount,
			Rate:       rate,
		},
		SourceAccountID: sourceAccountID,
		TargetAccountID: targetAccountID,
	}, nil
}

func FXConvertCash(w *world.World, humanID string, bankID string, fromCurrency string, toCurrency string, amount int) (FXCashConversionResult, error) {
	fromCurrency = strings.ToUpper(fromCurrency)
	toCurrency = strings.ToUpper(toCurrency)

	if amount <= 0 {
		return FXCashConversionResult{}, fmt.Errorf("amount must be greater than 0")
	}

	human, ok := w.Humans[humanID]
	if !ok {
		return FXCashConversionResult{}, fmt.Errorf("human does not exist: %s", humanID)
	}

	bank, ok := w.Banks[bankID]
	if !ok {
		return FXCashConversionResult{}, fmt.Errorf("bank does not exist: %s", bankID)
	}

	rate, convertedAmount, err := convertAmount(w, fromCurrency, toCurrency, amount)
	if err != nil {
		return FXCashConversionResult{}, err
	}

	if human.CashWallet[fromCurrency] < amount {
		return FXCashConversionResult{}, fmt.Errorf(
			"not enough cash in %s's wallet. Available: %d %s",
			humanID,
			human.CashWallet[fromCurrency],
			fromCurrency,
		)
	}

	if bank.CashVault[toCurrency] < convertedAmount {
		return FXCashConversionResult{}, fmt.Errorf(
			"not enough physical cash in %s's vault. Available: %d %s",
			bankID,
			bank.CashVault[toCurrency],
			toCurrency,
		)
	}

	human.CashWallet[fromCurrency] -= amount
	human.CashWallet[toCurrency] += convertedAmount
	bank.CashVault[fromCurrency] += amount
	bank.CashVault[toCurrency] -= convertedAmount

	return FXCashConversionResult{
		FXConversionResult: FXConversionResult{
			FromAmount: amount,
			ToAmount:   convertedAmount,
			Rate:       rate,
		},
	}, nil
}

func FXBankTrade(w *world.World, fromBankID string, toBankID string, fromCurrency string, toCurrency string, amount int) (FXBankTradeResult, error) {
	fromCurrency = strings.ToUpper(fromCurrency)
	toCurrency = strings.ToUpper(toCurrency)

	if amount <= 0 {
		return FXBankTradeResult{}, fmt.Errorf("amount must be greater than 0")
	}

	if fromBankID == toBankID {
		return FXBankTradeResult{}, fmt.Errorf("from bank and to bank must be different")
	}

	fromBank, ok := w.Banks[fromBankID]
	if !ok {
		return FXBankTradeResult{}, fmt.Errorf("from bank does not exist: %s", fromBankID)
	}

	toBank, ok := w.Banks[toBankID]
	if !ok {
		return FXBankTradeResult{}, fmt.Errorf("to bank does not exist: %s", toBankID)
	}

	rate, convertedAmount, err := convertAmount(w, fromCurrency, toCurrency, amount)
	if err != nil {
		return FXBankTradeResult{}, err
	}

	fromCentralBankID, fromCentralBank, err := centralBankForCurrency(w, fromCurrency)
	if err != nil {
		return FXBankTradeResult{}, err
	}

	toCentralBankID, toCentralBank, err := centralBankForCurrency(w, toCurrency)
	if err != nil {
		return FXBankTradeResult{}, err
	}

	fromBankFromReserves, err := reserveBalance(fromCentralBank, fromBank, fromCentralBankID, fromBankID)
	if err != nil {
		return FXBankTradeResult{}, err
	}

	if _, err := reserveBalance(fromCentralBank, toBank, fromCentralBankID, toBankID); err != nil {
		return FXBankTradeResult{}, err
	}

	toBankToReserves, err := reserveBalance(toCentralBank, toBank, toCentralBankID, toBankID)
	if err != nil {
		return FXBankTradeResult{}, err
	}

	if _, err := reserveBalance(toCentralBank, fromBank, toCentralBankID, fromBankID); err != nil {
		return FXBankTradeResult{}, err
	}

	if fromBankFromReserves < amount {
		return FXBankTradeResult{}, fmt.Errorf(
			"not enough %s reserves for %s at %s. Available: %d %s",
			fromCurrency,
			fromBankID,
			fromCentralBankID,
			fromBankFromReserves,
			fromCurrency,
		)
	}

	if toBankToReserves < convertedAmount {
		return FXBankTradeResult{}, fmt.Errorf(
			"not enough %s reserves for %s at %s. Available: %d %s",
			toCurrency,
			toBankID,
			toCentralBankID,
			toBankToReserves,
			toCurrency,
		)
	}

	fromCentralBank.ReserveAccounts[fromBankID] -= amount
	fromBank.ReserveBalances[fromCentralBankID] -= amount
	fromCentralBank.ReserveAccounts[toBankID] += amount
	toBank.ReserveBalances[fromCentralBankID] += amount

	toCentralBank.ReserveAccounts[toBankID] -= convertedAmount
	toBank.ReserveBalances[toCentralBankID] -= convertedAmount
	toCentralBank.ReserveAccounts[fromBankID] += convertedAmount
	fromBank.ReserveBalances[toCentralBankID] += convertedAmount

	return FXBankTradeResult{
		FXConversionResult: FXConversionResult{
			FromAmount: amount,
			ToAmount:   convertedAmount,
			Rate:       rate,
		},
		FromCentralBankID: fromCentralBankID,
		ToCentralBankID:   toCentralBankID,
	}, nil
}

func FXMarketID(fromCurrency string, toCurrency string) string {
	return strings.ToUpper(fromCurrency) + "_" + strings.ToUpper(toCurrency)
}

func convertAmount(w *world.World, fromCurrency string, toCurrency string, amount int) (float64, int, error) {
	rate, err := fxRate(w, fromCurrency, toCurrency)
	if err != nil {
		return 0, 0, err
	}

	return rate, int(math.Round(float64(amount) * rate)), nil
}

func fxRate(w *world.World, fromCurrency string, toCurrency string) (float64, error) {
	fromCurrency = strings.ToUpper(fromCurrency)
	toCurrency = strings.ToUpper(toCurrency)

	if fromCurrency == toCurrency {
		return 1, nil
	}

	if market, ok := w.FXMarkets[FXMarketID(fromCurrency, toCurrency)]; ok {
		return market.Rate, nil
	}

	if market, ok := w.FXMarkets[FXMarketID(toCurrency, fromCurrency)]; ok {
		return 1 / market.Rate, nil
	}

	return 0, fmt.Errorf("FX market does not exist: %s", FXMarketID(fromCurrency, toCurrency))
}

func centralBankForCurrency(w *world.World, currency string) (string, *world.CentralBank, error) {
	for centralBankID, centralBank := range w.CentralBanks {
		if centralBank.Currency == currency {
			return centralBankID, centralBank, nil
		}
	}

	return "", nil, fmt.Errorf("no central bank for %s", currency)
}
