package world

func New() *World {
	return &World{
		Version: 1,

		CentralBanks: map[string]*CentralBank{},
		Banks:        map[string]*Bank{},
		Humans:       map[string]*Human{},
		Accounts:     map[string]*Account{},
		Currencies:   map[string]*Currency{},
		Assets:       map[string]*Asset{},
		ReserveLoans: map[string]*ReserveLoan{},

		CustomerLoans: map[string]*CustomerLoan{},
		Holds:         map[string]map[string]any{},
		LedgerEntries: []map[string]any{},

		PaymentInstructions: map[string]*PaymentInstruction{},
		Messages:            map[string]map[string]any{},
		Settlements:         map[string]map[string]any{},

		Step2Systems:          map[string]map[string]any{},
		CorrespondentAccounts: map[string]map[string]any{},
		Bonds:                 map[string]map[string]any{},
		Cheques:               map[string]map[string]any{},
		CardAuthorizations:    map[string]map[string]any{},
		FXMarkets:             map[string]map[string]any{},
		Snapshots:             map[string]map[string]any{},

		CommandHistory: []CommandHistoryEntry{},
	}
}

func NewCurrency() *Currency {
	return &Currency{}
}

func NewAsset(id string, ownerType AssetOwnerType, ownerID string, currency string, estimatedValue int) *Asset {
	return &Asset{
		ID:                          id,
		OwnerType:                   ownerType,
		OwnerID:                     ownerID,
		Currency:                    currency,
		EstimatedValue:              estimatedValue,
		PledgedToCentralBankID:      "",
		CollateralForReserveLoanID:  "",
		PledgedToBankID:             "",
		CollateralForCustomerLoanID: "",
	}
}

func NewReserveLoan(id string, centralBankID string, bankID string, currency string, amount int, collateralAssetID string) *ReserveLoan {
	return &ReserveLoan{
		ID:                id,
		CentralBankID:     centralBankID,
		BankID:            bankID,
		Currency:          currency,
		Principal:         amount,
		Outstanding:       amount,
		CollateralAssetID: collateralAssetID,
		Status:            ReserveLoanOpen,
	}
}

func NewPaymentInstruction(id string, paymentType PaymentType, senderHumanID string, senderBankID string, recipientHumanID string, recipientBankID string, senderAccountID string, recipientAccountID string, centralBankID string, currency string, amount int) *PaymentInstruction {
	return &PaymentInstruction{
		ID:                 id,
		Type:               paymentType,
		Status:             PaymentCompleted,
		SenderHumanID:      senderHumanID,
		RecipientHumanID:   recipientHumanID,
		SenderBankID:       senderBankID,
		RecipientBankID:    recipientBankID,
		SenderAccountID:    senderAccountID,
		RecipientAccountID: recipientAccountID,
		CentralBankID:      centralBankID,
		Currency:           currency,
		Amount:             amount,
	}
}

func NewCustomerLoan(id string, bankID string, borrowerHumanID string, currency string, amount int) *CustomerLoan {
	return &CustomerLoan{
		ID:                   id,
		BankID:               bankID,
		BorrowerHumanID:      borrowerHumanID,
		Currency:             currency,
		OriginalPrincipal:    amount,
		OutstandingPrincipal: amount,
		OutstandingInterest:  0,
		TotalInterestAccrued: 0,
		WrittenOffPrincipal:  0,
		WrittenOffInterest:   0,
		CollateralAssetID:    "",
		Status:               CustomerLoanActive,
	}
}

func NewCentralBank(id string, currency string) *CentralBank {
	return &CentralBank{
		ID:              id,
		Name:            id,
		Currency:        currency,
		CashIssued:      0,
		CashVault:       0,
		ReserveAccounts: map[string]int{},
		LoansToBanks:    map[string]int{},
		Securities:      map[string]int{},
	}
}

func NewBank(id string) *Bank {
	return &Bank{
		ID:                    id,
		Name:                  id,
		CashVault:             map[string]int{},
		ReserveBalances:       map[string]int{},
		LoansFromCentralBanks: map[string]int{},
		CustomerAccounts:      []string{},
		CustomerLoans:         []string{},
		InterestIncome:        map[string]int{},
		LoanLossExpense:       map[string]int{},
		Equity:                map[string]int{},
		FXInventory:           map[string]int{},
	}
}

func NewHuman(id string) *Human {
	return &Human{
		ID:           id,
		Name:         id,
		CashWallet:   map[string]int{},
		BankAccounts: []string{},
		Loans:        map[string]int{},
	}
}

func NewAccount(id string, humanID string, bankID string, currency string) *Account {
	return &Account{
		ID:            id,
		OwnerHumanID:  humanID,
		BankID:        bankID,
		Currency:      currency,
		BookedBalance: 0,
		Holds:         []string{},
		Status:        AccountActive,
	}
}
