package kafkaaggregator

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/JakubDaleki/transfer-app/shared-dependencies"
	"github.com/segmentio/kafka-go"
)

type Aggregator struct {
	brokers    []string
	partitions []kafka.PartitionOffsets
}

func NewAggregator(brokers []string) *Aggregator {
	agg := new(Aggregator)
	agg.brokers = brokers
	addr := brokers[0]
	topic := "transfers"
	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		panic(err.Error())
	}

	client := kafka.Client{
		Addr:    kafka.TCP(brokers...),
		Timeout: 10 * time.Second,
	}

	partitionIds := []kafka.OffsetRequest{}
	for _, partition := range partitions {
		partitionIds = append(partitionIds, kafka.OffsetRequest{Partition: partition.ID, Timestamp: kafka.LastOffset})
	}
	offRequest := map[string][]kafka.OffsetRequest{topic: partitionIds}
	r, err := client.ListOffsets(context.Background(), &kafka.ListOffsetsRequest{Topics: offRequest})
	if err != nil {
		panic(err.Error())
	}
	agg.partitions = r.Topics[topic]
	return agg

}

func (agg *Aggregator) ProcessAllPartitions(ch chan map[string]float64) {
	log.Print("Processing data from ", len(agg.partitions), " partitions")
	for _, partition := range agg.partitions {
		balances := agg.ProcessPartition(partition.Partition, partition.LastOffset)
		ch <- balances
	}

	close(ch)
	log.Print("Successfully processed data from ", len(agg.partitions), " partitions")

}

func (agg *Aggregator) ProcessPartition(partitionId int, offset int64) map[string]float64 {
	aggregations := make(map[string]float64)
	kafkaR := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   agg.brokers,
		Topic:     "transfers",
		Partition: partitionId,
	})
	kafkaR.SetOffset(0)
	for i := int64(0); i < offset; i++ {
		m, _ := kafkaR.FetchMessage(context.Background())
		mKey := string(m.Key)
		transfer := new(shared.Transfer)
		json.Unmarshal(m.Value, transfer)
		currAmount := aggregations[mKey]

		if mKey == transfer.From {
			aggregations[mKey] = currAmount - transfer.Amount

		} else {
			aggregations[mKey] = currAmount + transfer.Amount
		}
	}

	return aggregations
}
