package shared

type Transfer struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type Balance struct {
	Username string `json:"username"`
	Balance  int    `json:"balance"`
}
