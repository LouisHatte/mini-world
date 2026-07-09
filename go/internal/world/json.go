package world

import "encoding/json"

func (w *World) UnmarshalJSON(data []byte) error {
	type alias World

	defaults := New()
	if err := json.Unmarshal(data, (*alias)(defaults)); err != nil {
		return err
	}

	*w = *defaults
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
