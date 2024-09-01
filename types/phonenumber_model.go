package types

type PhoneNumberModel struct {
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Number      string `json:"number"`
}
