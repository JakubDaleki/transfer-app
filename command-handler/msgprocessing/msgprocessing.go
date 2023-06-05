package msgprocessing

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/JakubDaleki/transfer-app/shared-dependencies"
	pb "github.com/JakubDaleki/transfer-app/shared-dependencies/grpc"
	"github.com/segmentio/kafka-go"
)

type TransferProcessor struct {
	Client pb.QueryServiceClient
	KafkaR *kafka.Reader
	KafkaW *kafka.Writer
}

func (t *TransferProcessor) Process() {
	for {
		m, err := t.KafkaR.FetchMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}
		mKey := string(m.Key)
		transfer := new(shared.Transfer)
		json.Unmarshal(m.Value, transfer)
		if mKey == transfer.From {
			// we are performing subtraction
			_, err = t.Client.UpdateBalance(context.Background(), &pb.UpdateBalanceRequest{User: transfer.From, Amount: -transfer.Amount})
			if err != nil {
				fmt.Printf("Failed to decrease balance of %v by %v. Reason: %v\n", transfer.From, transfer.Amount, err)
			} else {
				fmt.Printf("Subtracted %v from %v's account.\n", transfer.Amount, transfer.From)
				// we need to send a copy to increase other user balance and store it on their partition
				_ = t.KafkaW.WriteMessages(context.Background(), kafka.Message{
					Value: m.Value,
					Key:   []byte(transfer.To),
				})
			}

		} else {
			// we are performing addition
			_, err = t.Client.UpdateBalance(context.Background(), &pb.UpdateBalanceRequest{User: transfer.To, Amount: transfer.Amount})
			if err != nil {
				fmt.Printf("Failed to increase balance of %v by %v. Reason: %v\n", transfer.To, transfer.Amount, err)
			} else {
				fmt.Printf("transfered %v from %v to %v\n", transfer.Amount, transfer.From, transfer.To)
			}
		}

		t.KafkaR.CommitMessages(context.Background(), m)

	}
}
