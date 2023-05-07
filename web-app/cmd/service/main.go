package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JakubDaleki/transfer-app/webapp/api/handlers"
	"github.com/JakubDaleki/transfer-app/webapp/api/middleware"
	"github.com/JakubDaleki/transfer-app/webapp/utils/db"
	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
)

func main() {
	counter := 0

	// wait for kafka service to be up
	conn, err := kafka.Dial("tcp", "broker:29092")
	for err != nil {
		log.Println("failed to dial leader:", err)
		time.Sleep(time.Second * 10)
		conn, err = kafka.Dial("tcp", "broker:29092")
		counter++
		if counter == 4 {
			return
		}
	}
	// we can close it as we are going to use high-level Writer API
	conn.Close()

	urlExample := "postgres://postgres:password123@db:5432/postgres"
	conn2, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	var name string
	var weight int64
	_, err = conn2.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS users (id uuid NOT NULL, username text UNIQUE NOT NULL, password text NOT NULL);")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(name, weight)

	// round-robin writer
	kafkaW := &kafka.Writer{
		Addr:  kafka.TCP("broker:29092"),
		Topic: "transfers",
	}
	connector := db.NewConnector()

	// use container DI like uber-go/dig instead of manual
	DIBalanceHandler := func(w http.ResponseWriter, r *http.Request, username string) {
		handlers.BalanceHandler(w, r, username, connector)
	}
	DIAuthHandler := func(w http.ResponseWriter, r *http.Request) { handlers.AuthHandler(w, r, connector) }
	DIRegHandler := func(w http.ResponseWriter, r *http.Request) { handlers.RegHandler(w, r, connector) }
	DITransferHandler := func(w http.ResponseWriter, r *http.Request, username string) {
		handlers.TransferHandler(w, r, username, kafkaW)
	}

	http.HandleFunc("/balance", middleware.AuthMiddleware(DIBalanceHandler))
	http.HandleFunc("/authentication", DIAuthHandler)
	http.HandleFunc("/register", DIRegHandler)
	http.HandleFunc("/transfer", middleware.AuthMiddleware(DITransferHandler))
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))

}
