package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var connector = NewConnector()

func main() {
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
	err := connector.TransferTx(username, transfer.To, transfer.Amount)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	w.Write([]byte(fmt.Sprintf("{\"message\": \"Transfered %v from %v to %v\"}", transfer.Amount, username, transfer.To)))
}
