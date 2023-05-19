package shared

import (
	"github.com/google/uuid"
)

type Transfer struct {
	Id     uuid.UUID `json:"id"`
	From   string    `json:"from"`
	To     string    `json:"to"`
	Amount int       `json:"amount"`
	Status string    `json:"status"`
}

type Balance struct {
	Username string `json:"username"`
	Balance  int    `json:"balance"`
}
