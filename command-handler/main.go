package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/JakubDaleki/transfer-app/shared-dependencies"
	pb "github.com/JakubDaleki/transfer-app/shared-dependencies/grpc"
	kafkautils "github.com/JakubDaleki/transfer-app/shared-dependencies/kafka"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func transferProcessing(client pb.QueryServiceClient, kafkaR *kafka.Reader, kafkaW *kafka.Writer) {
	for {
		m, err := kafkaR.FetchMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}
		mKey := string(m.Key)
		transfer := new(shared.Transfer)
		json.Unmarshal(m.Value, transfer)
		if mKey == transfer.From {
			// we are performing subtraction
			_, err = client.UpdateBalance(context.Background(), &pb.UpdateBalanceRequest{User: transfer.From, Amount: -transfer.Amount})
			if err != nil {
				fmt.Printf("Failed to decrease balance of %v by %v. Reason: %v\n", transfer.From, transfer.Amount, err)
			} else {
				fmt.Printf("Subtracted %v from %v's account.\n", transfer.Amount, transfer.From)
				// we need to send a copy to increase other user balance and store it on their partition
				_ = kafkaW.WriteMessages(context.Background(), kafka.Message{
					Value: m.Value,
					Key:   []byte(transfer.To),
				})
			}

		} else {
			// we are performing addition
			_, err = client.UpdateBalance(context.Background(), &pb.UpdateBalanceRequest{User: transfer.To, Amount: transfer.Amount})
			if err != nil {
				fmt.Printf("Failed to increase balance of %v by %v. Reason: %v\n", transfer.To, transfer.Amount, err)
			} else {
				fmt.Printf("transfered %v from %v to %v\n", transfer.Amount, transfer.From, transfer.To)
			}
		}

		kafkaR.CommitMessages(context.Background(), m)

	}
}

func main() {
	conn, err := grpc.Dial("queryservice:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	err = kafkautils.WaitForKafka()

	if err != nil {
		panic(err.Error())
	}

	kafkaR := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"broker:29092"},
		GroupID: "transfer-processors-group",
		Topic:   "transfers",
	})

	defer kafkaR.Close()

	kafkaW := &kafka.Writer{
		Addr:     kafka.TCP("broker:29092"),
		Topic:    "transfers",
		Balancer: &kafka.Hash{},
	}

	defer kafkaW.Close()

	client := pb.NewQueryServiceClient(conn)
	transferProcessing(client, kafkaR, kafkaW)
}
