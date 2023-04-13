package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// this should be a database instead
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var userAmounts = map[string]int{
	"user1": 0,
	"user2": 0,
}

func main() {
	http.HandleFunc("/balance", authMiddleware(balanceHandler))
	http.HandleFunc("/authentication", authHandler)
	http.HandleFunc("/transfer", authMiddleware(transferHandler))
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))

}

func balanceHandler(w http.ResponseWriter, r *http.Request, username string) {
	balance := userAmounts[username]
	w.Write([]byte(fmt.Sprintf("{\"balance\": \"%v\"}", balance)))
}

type Transfer struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

func transferHandler(w http.ResponseWriter, r *http.Request, username string) {
	transfer := new(Transfer)
	json.NewDecoder(r.Body).Decode(transfer)
	userAmounts[username] -= transfer.Amount
	userAmounts[transfer.To] -= transfer.Amount
	w.Write([]byte(fmt.Sprintf("{\"balance\": \"%v\"}", userAmounts[username])))
}
