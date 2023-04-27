package main

import (
	"context"
	"errors"
	"github.com/JakubDaleki/transfer-app/webapp/api/handlers"
	"github.com/JakubDaleki/transfer-app/webapp/api/middleware"
	"github.com/JakubDaleki/transfer-app/webapp/utils/db"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"time"
)

var topic = "my-topic"
var partition = 0
var connector = db.NewConnector()
var conn *kafka.Conn

func DIWrapper() func(http.ResponseWriter, *http.Request, string) {
	return func(w http.ResponseWriter, r *http.Request, username string) {
		handlers.BalanceHandler(w, r, username, connector)
	}
}

func main() {
	counter := 0
	var err = errors.New("Non nil error")
	for err != nil {
		log.Println("failed to dial leader:", err)
		time.Sleep(time.Second * 10)
		conn, err = kafka.DialLeader(context.Background(), "tcp", "broker:29092", topic, partition)
		counter++
		if counter == 4 {
			return
		}
	}

	defer conn.Close()

	// use container DI like uber-go/dig instead of manual
	DIBalanceHandler := func(w http.ResponseWriter, r *http.Request, username string) {
		handlers.BalanceHandler(w, r, username, connector)
	}
	DIAuthHandler := func(w http.ResponseWriter, r *http.Request) { handlers.AuthHandler(w, r, connector) }
	DIRegHandler := func(w http.ResponseWriter, r *http.Request) { handlers.RegHandler(w, r, connector) }
	DITransferHandler := func(w http.ResponseWriter, r *http.Request, username string) {
		handlers.TransferHandler(w, r, username, connector, conn)
	}

	http.HandleFunc("/balance", middleware.AuthMiddleware(DIBalanceHandler))
	http.HandleFunc("/authentication", DIAuthHandler)
	http.HandleFunc("/register", DIRegHandler)
	http.HandleFunc("/transfer", middleware.AuthMiddleware(DITransferHandler))
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))

}
