package dto

type BinCardResponse struct {
	Brand   string `json:"brand"`
	Type    string `json:"type"`
	Country struct {
		Name string `json:"name"`
	} `json:"country"`
	Bank struct {
		Name string `json:"name"`
	} `json:"bank"`
}
