package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var connector = NewConnector()

func main() {
	http.HandleFunc("/balance", balanceHandler)
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))

}

func balanceHandler(w http.ResponseWriter, r *http.Request, username string) {
	balance := connector.userAmounts[username]
	w.Write([]byte(fmt.Sprintf("{\"balance\": \"%v\"}", balance)))
}
