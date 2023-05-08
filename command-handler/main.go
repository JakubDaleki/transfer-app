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

func updateBalance(pool *pgxpool.Pool, username string) error {
	_, err := pool.Exec(context.Background(), "update balance set balance=$1 where id=$2", description, itemNum)
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
		json.NewDecoder(string(m.Value)).Decode(transfer)
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func main() {
	url := "postgres://postgres:password123@db:5432/postgres"
	conn, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return
	}

}
