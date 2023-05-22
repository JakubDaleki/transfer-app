package shared

import (
	"github.com/google/uuid"
)

type Transfer struct {
	Id     uuid.UUID `json:"id"`
	From   string    `json:"from"`
	To     string    `json:"to"`
	Amount float64   `json:"amount"`
	Status string    `json:"status"`
}

type Balance struct {
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
}
