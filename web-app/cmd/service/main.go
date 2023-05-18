package main

import (
	"log"
	"net/http"

	"github.com/JakubDaleki/transfer-app/webapp/api/router"
	"github.com/JakubDaleki/transfer-app/webapp/utils"
	"github.com/JakubDaleki/transfer-app/webapp/utils/db"
	"github.com/segmentio/kafka-go"
)

func main() {
	// wait for kafka service to be up
	err := utils.WaitForKafka()
	if err != nil {
		panic(err)
	}

	connector, err := db.WaitForDb()

	if err != nil {
		panic(err)
	}

	// round-robin writer
	kafkaW := &kafka.Writer{
		Addr:  kafka.TCP("broker:29092"),
		Topic: "transfers",
	}

	s := &http.Server{
		Addr:    ":8000",
		Handler: router.New(connector, kafkaW),
	}

	// start the server on port 8000
	log.Fatal(s.ListenAndServe())

}
