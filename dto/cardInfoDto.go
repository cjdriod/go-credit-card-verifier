package dto

type CardInfoDto struct {
	Type          string   `json:"type"`
	Issuer        string   `json:"issuer"`
	IssueCountry  string   `json:"issueCountry"`
	Network       string   `json:"network"`
	CardNumber    string   `json:"cardNumber"`
	Activities    []string `json:"activities"`
	IsBlackListed bool     `json:"isBlackListed"`
}
