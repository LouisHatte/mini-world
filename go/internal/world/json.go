package world

import "encoding/json"

func (w *World) UnmarshalJSON(data []byte) error {
	type alias World

	defaults := New()
	if err := json.Unmarshal(data, (*alias)(defaults)); err != nil {
		return err
	}

	if defaults.Currencies == nil {
		defaults.Currencies = map[string]*Currency{}
	}

	for _, centralBank := range defaults.CentralBanks {
		if centralBank.Currency != "" {
			if _, exists := defaults.Currencies[centralBank.Currency]; !exists {
				defaults.Currencies[centralBank.Currency] = NewCurrency()
			}
		}
	}

	*w = *defaults
	return nil
}

func (currency *Currency) UnmarshalJSON(data []byte) error {
	type alias Currency

	defaults := NewCurrency()
	if err := json.Unmarshal(data, (*alias)(defaults)); err != nil {
		return err
	}

	*currency = *defaults
	return nil
}

func (asset *Asset) UnmarshalJSON(data []byte) error {
	type alias Asset

	defaults := NewAsset("", "", "", "", 0)
	if err := json.Unmarshal(data, (*alias)(defaults)); err != nil {
		return err
	}

	*asset = *defaults
	return nil
}

func (reserveLoan *ReserveLoan) UnmarshalJSON(data []byte) error {
	type alias ReserveLoan

	defaults := NewReserveLoan("", "", "", "", 0, "")
	if err := json.Unmarshal(data, (*alias)(defaults)); err != nil {
		return err
	}

	*reserveLoan = *defaults
	return nil
}

func (paymentInstruction *PaymentInstruction) UnmarshalJSON(data []byte) error {
	type alias PaymentInstruction

	defaults := NewPaymentInstruction("", "", "", "", "", "", "", "", "", "", 0)
	if err := json.Unmarshal(data, (*alias)(defaults)); err != nil {
		return err
	}

	*paymentInstruction = *defaults
	return nil
}

func (customerLoan *CustomerLoan) UnmarshalJSON(data []byte) error {
	type alias CustomerLoan

	defaults := NewCustomerLoan("", "", "", "", 0)
	if err := json.Unmarshal(data, (*alias)(defaults)); err != nil {
		return err
	}

	*customerLoan = *defaults
	return nil
}

func (fxMarket *FXMarket) UnmarshalJSON(data []byte) error {
	type alias FXMarket

	defaults := NewFXMarket("", "", "", 0)
	if err := json.Unmarshal(data, (*alias)(defaults)); err != nil {
		return err
	}

	*fxMarket = *defaults
	return nil
}

func (correspondentAccount *CorrespondentAccount) UnmarshalJSON(data []byte) error {
	type alias CorrespondentAccount

	defaults := NewCorrespondentAccount("", "", "", "")
	if err := json.Unmarshal(data, (*alias)(defaults)); err != nil {
		return err
	}

	*correspondentAccount = *defaults
	return nil
}

func (centralBank *CentralBank) UnmarshalJSON(data []byte) error {
	type alias CentralBank

	defaults := NewCentralBank("", "")
	if err := json.Unmarshal(data, (*alias)(defaults)); err != nil {
		return err
	}

	*centralBank = *defaults
	return nil
}

func (bank *Bank) UnmarshalJSON(data []byte) error {
	type alias Bank

	defaults := NewBank("")
	if err := json.Unmarshal(data, (*alias)(defaults)); err != nil {
		return err
	}

	*bank = *defaults
	return nil
}

func (human *Human) UnmarshalJSON(data []byte) error {
	type alias Human

	defaults := NewHuman("")
	if err := json.Unmarshal(data, (*alias)(defaults)); err != nil {
		return err
	}

	*human = *defaults
	return nil
}

func (account *Account) UnmarshalJSON(data []byte) error {
	type alias Account

	defaults := NewAccount("", "", "", "")
	if err := json.Unmarshal(data, (*alias)(defaults)); err != nil {
		return err
	}

	*account = *defaults
	return nil
}
