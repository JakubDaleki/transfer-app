package db

import (
	"fmt"

	"github.com/JakubDaleki/transfer-app/shared-dependencies"
	"github.com/hashicorp/go-memdb"
)

type Database struct {
	memdb *memdb.MemDB
}

func (db *Database) GetBalance(username string) *shared.Balance {
	txn := db.memdb.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("balance", "id", username)
	if err != nil || raw == nil {
		return &shared.Balance{Username: username, Balance: 0}
	}

	return raw.(*shared.Balance)
}

func (db *Database) UpdateBalance(balanceChange shared.Balance) error {
	txn := db.memdb.Txn(true)
	balance, _ := txn.First("balance", "id", balanceChange.Username)
	var newBalance float64
	if balance == nil {
		newBalance = balanceChange.Balance
	} else {
		newBalance = balance.(*shared.Balance).Balance + balanceChange.Balance
	}
	if newBalance < 0 {
		txn.Abort()
		return fmt.Errorf("not enough funds")
	}

	if err := txn.Insert("balance", &shared.Balance{Username: balanceChange.Username, Balance: newBalance}); err != nil {
		txn.Abort()
		return fmt.Errorf(fmt.Sprintf("couldn't update %s balance", balanceChange.Username))
	}

	txn.Commit()
	return nil
}

func NewDatabase() (*Database, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"balance": &memdb.TableSchema{
				Name: "balance",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Username"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}
	txn := db.Txn(true)
	txn.Insert("balance", &shared.Balance{Username: "user1", Balance: 528.11})
	txn.Commit()
	return &Database{memdb: db}, nil
}
