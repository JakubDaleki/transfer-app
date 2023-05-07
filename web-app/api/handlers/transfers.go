package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JakubDaleki/transfer-app/webapp/api/resource/auth"
	"github.com/JakubDaleki/transfer-app/webapp/api/resource/transfers"
	"github.com/JakubDaleki/transfer-app/webapp/utils/db"
	"github.com/segmentio/kafka-go"
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

func TransferHandler(w http.ResponseWriter, r *http.Request, username string, kafkaW *kafka.Writer) {
	if r.Method == "POST" {
		transfer := new(transfers.Transfer)
		json.NewDecoder(r.Body).Decode(transfer)
		txByteMsg, err := json.Marshal(*transfer)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
			return
		}

		err = kafkaW.WriteMessages(context.Background(), kafka.Message{Value: txByteMsg})

		if err != nil {
			w.WriteHeader(http.StatusFailedDependency)
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("{\"message\": \"Queued transfer of %v from %v to %v\"}", transfer.Amount, username, transfer.To)))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("{\"message\": \"This HTTP method not allowed, use POST instead.")))
	}

}
