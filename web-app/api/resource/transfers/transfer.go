package transfers

type Transfer struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}
