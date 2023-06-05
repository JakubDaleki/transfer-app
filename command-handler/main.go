package main

import (
	"log"

	"github.com/JakubDaleki/transfer-app/command-handler/msgprocessing"
	pb "github.com/JakubDaleki/transfer-app/shared-dependencies/grpc"
	kafkautils "github.com/JakubDaleki/transfer-app/shared-dependencies/kafka"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

	brokers := kafkautils.GetBootstrapServers()
	kafkaR := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: "transfer-processors-group",
		Topic:   "transfers",
	})

	defer kafkaR.Close()

	kafkaW := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    "transfers",
		Balancer: &kafka.Hash{},
	}

	defer kafkaW.Close()

	client := pb.NewQueryServiceClient(conn)

	proc := msgprocessing.TransferProcessor{
		Client: client,
		KafkaR: kafkaR,
		KafkaW: kafkaW,
	}
	proc.Process()
}
