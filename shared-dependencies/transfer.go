package shared

type Transfer struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}
