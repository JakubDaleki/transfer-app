package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"time"
)

var connector = NewConnector()
var conn *kafka.Conn

func main() {
	topic := "my-topic"
	partition := 0
	//var err error
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(1000 * time.Second))
	//defer conn.Close()

	http.HandleFunc("/balance", authMiddleware(balanceHandler))
	http.HandleFunc("/authentication", authHandler)
	http.HandleFunc("/register", regHandler)
	http.HandleFunc("/transfer", authMiddleware(transferHandler))
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))

}

func balanceHandler(w http.ResponseWriter, r *http.Request, username string) {
	balance := connector.userAmounts[username]
	w.Write([]byte(fmt.Sprintf("{\"balance\": \"%v\"}", balance)))
}

func regHandler(w http.ResponseWriter, r *http.Request) {
	cred := new(Credentials)
	json.NewDecoder(r.Body).Decode(cred)
	connector.AddNewUser(cred.Username, cred.Password)
	w.Write([]byte(fmt.Sprintf("{\"message\": \"ok\"}")))
}

type Transfer struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

func transferHandler(w http.ResponseWriter, r *http.Request, username string) {
	transfer := new(Transfer)
	json.NewDecoder(r.Body).Decode(transfer)
	/*
		err := connector.TransferTx(username, transfer.To, transfer.Amount)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
			return
		}
	*/
	// send transfer message to kafka
	// no partition id and no key for a message to go with round robin
	b, err := json.Marshal(*transfer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(b)
	_, err = conn.WriteMessages(
		kafka.Message{Value: b},
	)
	if err != nil {
		fmt.Println("failed to write messages:", err)
	}

	w.Write([]byte(fmt.Sprintf("{\"message\": \"Queued transfer of %v from %v to %v\"}", transfer.Amount, username, transfer.To)))
}
