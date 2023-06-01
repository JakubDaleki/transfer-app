package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/JakubDaleki/transfer-app/transfer-aggregator/kafkaaggregator"

	pb "github.com/JakubDaleki/transfer-app/shared-dependencies/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	mode := os.Args[1]
	if mode != "full" && mode != "runtime" {
		panic("Incorrent mode selected, please use 'full' or 'runtime'")
	}

	conn, err := grpc.Dial("queryservice:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewQueryServiceClient(conn)

	ch := make(chan map[string]float64)
	agg := new(kafkaaggregator.Aggregator)
	agg.ProcessAllPartitions(ch)

	for batchedBalance := range ch {
		_, err = client.RecreateBalances(context.Background(), &pb.BalancesMapRequest{BatchedBalances: batchedBalance})
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	fmt.Sprintln("Successfully processed data from ", len(agg.Partitions), " partitions")
}