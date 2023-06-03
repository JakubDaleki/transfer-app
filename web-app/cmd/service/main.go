package main

import (
	"log"
	"net/http"

	pb "github.com/JakubDaleki/transfer-app/shared-dependencies/grpc"
	kafkautils "github.com/JakubDaleki/transfer-app/shared-dependencies/kafka"
	"github.com/JakubDaleki/transfer-app/webapp/api/router"
	"github.com/JakubDaleki/transfer-app/webapp/utils/db"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// wait for kafka service to be up
	err := kafkautils.WaitForKafka()
	if err != nil {
		panic(err)
	}

	connector, err := db.WaitForDb()

	if err != nil {
		panic(err)
	}

	// use hash balancer to ensure partition and ordering per user
	kafkaW := &kafka.Writer{
		Addr:     kafka.TCP(kafkautils.GetBootstrapServers()...),
		Topic:    "transfers",
		Balancer: &kafka.Hash{},
	}

	conn, err := grpc.Dial("queryservice:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewQueryServiceClient(conn)

	s := &http.Server{
		Addr:    ":8000",
		Handler: router.New(connector, kafkaW, client),
	}

	// start the server on port 8000
	log.Fatal(s.ListenAndServe())

}
