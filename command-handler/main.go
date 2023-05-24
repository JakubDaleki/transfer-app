package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/JakubDaleki/transfer-app/shared-dependencies"
	pb "github.com/JakubDaleki/transfer-app/shared-dependencies/grpc"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func transferProcessing(client pb.QueryServiceClient) {
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
		transfer := new(shared.Transfer)
		json.Unmarshal(m.Value, transfer)
		_, err = client.MakeTransfer(context.Background(), &pb.TransferRequest{From: transfer.From, To: transfer.To, Amount: transfer.Amount})
		if err != nil {
			fmt.Printf("Failed to transfer %v from %v to %v\n", transfer.Amount, transfer.From, transfer.To)
			continue
		} else {
			fmt.Printf("transfered %v from %v to %v\n", transfer.Amount, transfer.From, transfer.To)
		}

	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func main() {
	conn, err := grpc.Dial("queryservice:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewQueryServiceClient(conn)
	transferProcessing(client)
}
