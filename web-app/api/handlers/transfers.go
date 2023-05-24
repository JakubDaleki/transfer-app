package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JakubDaleki/transfer-app/shared-dependencies"
	pb "github.com/JakubDaleki/transfer-app/shared-dependencies/grpc"
	"github.com/JakubDaleki/transfer-app/webapp/api/resource/auth"
	"github.com/JakubDaleki/transfer-app/webapp/utils/db"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func BalanceHandler(w http.ResponseWriter, r *http.Request, client pb.QueryServiceClient) {
	username := r.Context().Value("username").(string)
	response, _ := client.GetBalance(context.Background(), &pb.BalanceRequest{Username: username})
	json.NewEncoder(w).Encode(&shared.Balance{Username: username, Balance: response.Balance})
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
	transfer.Id = uuid.New()
	transfer.Status = "queued"
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
	json.NewEncoder(w).Encode(transfer)
}
