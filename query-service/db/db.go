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

	raw, err := txn.First("balance", "username", username)
	if err != nil {
		return &shared.Balance{Username: username, Balance: 0}
	}
	return raw.(*shared.Balance)
}

func (db *Database) MakeTransfer(transfer shared.Transfer) error {
	txn := db.memdb.Txn(true)
	rawUserFrom, err := txn.First("balance", "username", transfer.From)
	if err == memdb.ErrNotFound {
		txn.Abort()
		return fmt.Errorf("not enough funds")
	}
	newBalanceFrom := rawUserFrom.(*shared.Balance).Balance - transfer.Amount

	rawUserTo, err := txn.First("balance", "username", transfer.To)
	newBalanceTo := transfer.Amount
	if err == memdb.ErrNotFound {
		newBalanceTo += rawUserTo.(*shared.Balance).Balance
	}

	if err := txn.Insert("person", &shared.Balance{Username: transfer.From, Balance: newBalanceFrom}); err != nil {
		txn.Abort()
		return fmt.Errorf("couldn't update senders balance")
	}
	if err := txn.Insert("person", &shared.Balance{Username: transfer.To, Balance: newBalanceTo}); err != nil {
		txn.Abort()
		return fmt.Errorf("couldn't update receivers balance")
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
					"username": &memdb.IndexSchema{
						Name:    "username",
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

	return &Database{memdb: db}, nil
}
