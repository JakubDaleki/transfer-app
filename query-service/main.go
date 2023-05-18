package main

import (
	"github.com/JakubDaleki/transfer-app/query-service/db"
)

func main() {
	database, err := db.NewDatabase()
	if err != nil {
		panic(err)
	}

}
