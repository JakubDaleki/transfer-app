package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

// this should be a dependency
type Transfer struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

func updateBalance(pool *pgxpool.Pool, username string, amount int) error {
	tx, _ := pool.Begin(context.Background())
	tx.Exec(context.Background(), "update balance set balance=balance+$1 where username=$2", amount, username)
	err := tx.Commit(context.Background())
	return err
}

func transferProcessing(pool *pgxpool.Pool) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"broker:29092"},
		GroupID: "transfer-processors-group",
		Topic:   "transfers",
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		transfer := new(Transfer)
		json.Unmarshal(m.Value, transfer)
		err = updateBalance(pool, transfer.From, -transfer.Amount)
		if err != nil {
			fmt.Printf("Failed transfer of %v from %v to %v\n", transfer.Amount, transfer.From, transfer.To)
			continue
		}

		err = updateBalance(pool, transfer.To, transfer.Amount)
		fmt.Printf("transfered %v from %v to %v\n", transfer.Amount, transfer.From, transfer.To)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func main() {
	url := "postgres://postgres:password123@db:5432/postgres"
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return
	}
	transferProcessing(pool)
}
