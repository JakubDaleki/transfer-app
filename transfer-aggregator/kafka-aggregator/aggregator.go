package kafkaaggregator

import (
	"context"
	"encoding/json"

	"github.com/JakubDaleki/transfer-app/shared-dependencies"
	"github.com/segmentio/kafka-go"
)

type Aggregator struct {
	Partitions []kafka.Partition
}

func NewAggregator() *Aggregator {
	agg := new(Aggregator)

	conn, err := kafka.Dial("tcp", "broker:29092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions("transfers")
	if err != nil {
		panic(err.Error())
	}

	agg.Partitions = partitions
	return agg

}

func (agg *Aggregator) ProcessAllPartitions() {
	for _, partition := range agg.Partitions {
		agg.ProcessPartition(partition.ID)
	}

}

func (agg *Aggregator) ProcessPartition(partitionId int) map[string]float64 {
	aggregations := make(map[string]float64)
	kafkaR := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"broker:29092"},
		Topic:     "transfers",
		Partition: partitionId,
	})
	offset := kafkaR.Offset()
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
