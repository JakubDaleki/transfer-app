package main

import (
	"github.com/JakubDaleki/transfer-app/query-service/db"
)

func main() {
	_, err := db.NewDatabase()
	if err != nil {
		panic(err)
	}

}
