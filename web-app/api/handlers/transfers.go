package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/JakubDaleki/transfer-app/webapp/api/resource/auth"
	"github.com/JakubDaleki/transfer-app/webapp/api/resource/transfers"
	"github.com/JakubDaleki/transfer-app/webapp/utils/db"
	"github.com/segmentio/kafka-go"
	"net/http"
)

func BalanceHandler(w http.ResponseWriter, r *http.Request, username string, connector *db.Connector) {
	balance := connector.UserAmounts[username]
	w.Write([]byte(fmt.Sprintf("{\"balance\": \"%v\"}", balance)))
}

func RegHandler(w http.ResponseWriter, r *http.Request, connector *db.Connector) {
	cred := new(auth.Credentials)
	json.NewDecoder(r.Body).Decode(cred)
	connector.AddNewUser(cred.Username, cred.Password)
	w.Write([]byte(fmt.Sprintf("{\"message\": \"ok\"}")))
}

func TransferHandler(w http.ResponseWriter, r *http.Request, username string, connector *db.Connector, conn *kafka.Conn) {
	transfer := new(transfers.Transfer)
	json.NewDecoder(r.Body).Decode(transfer)
	err := connector.TransferTx(username, transfer.To, transfer.Amount)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}
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
