package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
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
		// todo process each message
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

}
