package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JakubDaleki/transfer-app/shared-dependencies"
	"github.com/JakubDaleki/transfer-app/webapp/api/resource/auth"
	"github.com/JakubDaleki/transfer-app/webapp/utils/db"
	"github.com/segmentio/kafka-go"
)

func BalanceHandler(w http.ResponseWriter, r *http.Request, connector *db.Connector) {
	username := r.Context().Value("username").(string)
	balance := connector.GetBalance(username)
	w.Write([]byte(fmt.Sprintf("{\"balance\": \"%v\"}", balance)))
}

func RegHandler(w http.ResponseWriter, r *http.Request, connector *db.Connector) {
	cred := new(auth.Credentials)
	json.NewDecoder(r.Body).Decode(cred)
	connector.AddNewUser(cred.Username, cred.Password)
	w.Write([]byte(fmt.Sprintf("{\"message\": \"ok\"}")))
}

func TransferHandler(w http.ResponseWriter, r *http.Request, kafkaW *kafka.Writer) {
	transfer := new(shared.Transfer)
	json.NewDecoder(r.Body).Decode(transfer)
	username := r.Context().Value("username").(string)
	transfer.From = username
	txByteMsg, err := json.Marshal(*transfer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	err = kafkaW.WriteMessages(context.Background(), kafka.Message{
		Value: txByteMsg,
		Key:   []byte(username),
	})

	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("{\"message\": \"Queued transfer of %v from %v to %v\"}", transfer.Amount, username, transfer.To)))

}
