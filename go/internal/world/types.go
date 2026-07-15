package world

type World struct {
	Version int `json:"version"`

	CentralBanks map[string]*CentralBank `json:"central_banks"`
	Banks        map[string]*Bank        `json:"banks"`
	Humans       map[string]*Human       `json:"humans"`
	Accounts     map[string]*Account     `json:"accounts"`
	Currencies   map[string]*Currency    `json:"currencies"`
	Assets       map[string]*Asset       `json:"assets"`
	ReserveLoans map[string]*ReserveLoan `json:"reserve_loans"`

	CustomerLoans map[string]map[string]any `json:"customer_loans"`
	Holds         map[string]map[string]any `json:"holds"`
	LedgerEntries []map[string]any          `json:"ledger_entries"`

	PaymentInstructions map[string]*PaymentInstruction `json:"payment_instructions"`
	Messages            map[string]map[string]any      `json:"messages"`
	Settlements         map[string]map[string]any      `json:"settlements"`

	Step2Systems          map[string]map[string]any `json:"step2_systems"`
	CorrespondentAccounts map[string]map[string]any `json:"correspondent_accounts"`
	Bonds                 map[string]map[string]any `json:"bonds"`
	Cheques               map[string]map[string]any `json:"cheques"`
	CardAuthorizations    map[string]map[string]any `json:"card_authorizations"`
	FXMarkets             map[string]map[string]any `json:"fx_markets"`
	Snapshots             map[string]map[string]any `json:"snapshots"`

	CommandHistory []CommandHistoryEntry `json:"command_history"`
}

type CentralBank struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Currency        string         `json:"currency"`
	CashIssued      int            `json:"cash_issued"`
	CashVault       int            `json:"cash_vault"`
	ReserveAccounts map[string]int `json:"reserve_accounts"`
	LoansToBanks    map[string]int `json:"loans_to_banks"`
	Securities      map[string]int `json:"securities"`
}

type Currency struct{}

type AssetOwnerType string

const (
	AssetOwnerHuman       AssetOwnerType = "HUMAN"
	AssetOwnerBank        AssetOwnerType = "BANK"
	AssetOwnerCentralBank AssetOwnerType = "CENTRAL_BANK"
)

type Asset struct {
	ID                         string         `json:"id"`
	OwnerType                  AssetOwnerType `json:"owner_type"`
	OwnerID                    string         `json:"owner_id"`
	Currency                   string         `json:"currency"`
	EstimatedValue             int            `json:"estimated_value"`
	PledgedToCentralBankID     string         `json:"pledged_to_central_bank_id"`
	CollateralForReserveLoanID string         `json:"collateral_for_reserve_loan_id"`
}

type ReserveLoanStatus string

const (
	ReserveLoanOpen   ReserveLoanStatus = "OPEN"
	ReserveLoanRepaid ReserveLoanStatus = "REPAID"
)

type ReserveLoan struct {
	ID                string            `json:"id"`
	CentralBankID     string            `json:"central_bank_id"`
	BankID            string            `json:"bank_id"`
	Currency          string            `json:"currency"`
	Principal         int               `json:"principal"`
	Outstanding       int               `json:"outstanding"`
	CollateralAssetID string            `json:"collateral_asset_id"`
	Status            ReserveLoanStatus `json:"status"`
}

type Bank struct {
	ID                    string         `json:"id"`
	Name                  string         `json:"name"`
	CashVault             map[string]int `json:"cash_vault"`
	ReserveBalances       map[string]int `json:"reserve_balances"`
	LoansFromCentralBanks map[string]int `json:"loans_from_central_banks"`
	CustomerAccounts      []string       `json:"customer_accounts"`
	CustomerLoans         []string       `json:"customer_loans"`
	InterestIncome        map[string]int `json:"interest_income"`
	LoanLossExpense       map[string]int `json:"loan_loss_expense"`
	Equity                map[string]int `json:"equity"`
	FXInventory           map[string]int `json:"fx_inventory"`
}

type Human struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	CashWallet   map[string]int `json:"cash_wallet"`
	BankAccounts []string       `json:"bank_accounts"`
	Loans        map[string]int `json:"loans"`
}

type AccountStatus string

const (
	AccountActive AccountStatus = "ACTIVE"
	AccountFrozen AccountStatus = "FROZEN"
	AccountClosed AccountStatus = "CLOSED"
)

type Account struct {
	ID            string        `json:"id"`
	OwnerHumanID  string        `json:"owner_human_id"`
	BankID        string        `json:"bank_id"`
	Currency      string        `json:"currency"`
	BookedBalance int           `json:"booked_balance"`
	Holds         []string      `json:"holds"`
	Status        AccountStatus `json:"status"`
}

type PaymentType string

const (
	PaymentInternal  PaymentType = "INTERNAL"
	PaymentInterbank PaymentType = "INTERBANK"
)

type PaymentStatus string

const (
	PaymentCompleted PaymentStatus = "COMPLETED"
)

type PaymentInstruction struct {
	ID                 string        `json:"id"`
	Type               PaymentType   `json:"type"`
	Status             PaymentStatus `json:"status"`
	SenderHumanID      string        `json:"sender_human_id"`
	RecipientHumanID   string        `json:"recipient_human_id"`
	SenderBankID       string        `json:"sender_bank_id"`
	RecipientBankID    string        `json:"recipient_bank_id"`
	SenderAccountID    string        `json:"sender_account_id"`
	RecipientAccountID string        `json:"recipient_account_id"`
	CentralBankID      string        `json:"central_bank_id"`
	Currency           string        `json:"currency"`
	Amount             int           `json:"amount"`
}

type CommandHistoryEntry struct {
	ID           int      `json:"id"`
	TimestampUTC string   `json:"timestamp_utc"`
	Command      *string  `json:"command"`
	Argv         []string `json:"argv"`
}
