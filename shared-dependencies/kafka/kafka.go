package utils

import (
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func WaitForKafka() error {
	for trial := 0; trial < 3; trial++ {
		conn, err := kafka.Dial("tcp", "broker:29092")
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(time.Second * 5)
	}

	return fmt.Errorf("couldn't connect to kafka")
}

func GetBootstrapServers() []string {
	// not all brokers have to be specified, only few
	return []string{"broker:29092", "broker2:29092", "broker3:29092"}
}
