package db

import (
	"testing"

	"github.com/JakubDaleki/transfer-app/shared-dependencies"
)

func TestGetBalance(t *testing.T) {
	database, _ := NewDatabase()
	balance := database.GetBalance("test")
	if balance.Balance != 0 {
		t.FailNow()
	}
}

func TestUpdateBalance(t *testing.T) {
	database, _ := NewDatabase()
	database.UpdateBalance(shared.Balance{Username: "test", Balance: 100})
	balance := database.GetBalance("test")
	if balance.Balance != 100 {
		t.FailNow()
	}
}
