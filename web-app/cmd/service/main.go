package main

import (
	"log"
	"net/http"
	"time"

	"github.com/JakubDaleki/transfer-app/webapp/api/router"
	"github.com/JakubDaleki/transfer-app/webapp/utils/db"
	"github.com/segmentio/kafka-go"
)

func main() {
	// wait for kafka service to be up
	conn, err := kafka.Dial("tcp", "broker:29092")
	for trial := 0; err != nil || trial == 3; trial++ {
		log.Println("failed to dial leader:", err)
		time.Sleep(time.Second * 10)
		conn, err = kafka.Dial("tcp", "broker:29092")
	}
	// we can close it as we are going to use high-level Writer API
	conn.Close()

	// round-robin writer
	kafkaW := &kafka.Writer{
		Addr:  kafka.TCP("broker:29092"),
		Topic: "transfers",
	}

	connector, err := db.NewConnector()
	for trial := 0; err != nil || trial == 3; trial++ {
		log.Println("failed to create db connector:", err)
		time.Sleep(time.Second * 10)
		connector, err = db.NewConnector()
	}

	s := &http.Server{
		Addr:    ":8000",
		Handler: router.New(connector, kafkaW),
	}

	// start the server on port 8000
	log.Fatal(s.ListenAndServe())

}
