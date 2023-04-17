package main

import (
	"errors"
	"sync"
)

// this should be a database instead
type Connector struct {
	users                    map[string]string
	userAmounts              map[string]int
	usersMutex, balanceMutex sync.Mutex
}

func (conn *Connector) GetPassword(username string) string {
	// what if user doesnt exist
	conn.usersMutex.Lock()
	defer conn.usersMutex.Unlock()
	return conn.users[username]
}

func (conn *Connector) TransferTx(from, to string, amount int) error {
	// what if user doesnt exist
	conn.balanceMutex.Lock()
	defer conn.balanceMutex.Unlock()
	if conn.userAmounts[from] < amount {
		return errors.New("Not enough funds on your account")
	}
	conn.userAmounts[from] -= amount
	conn.userAmounts[to] += amount
	return nil
}

func (conn *Connector) AddNewUser(username, password string) {
	conn.usersMutex.Lock()
	conn.users[username] = password
	conn.usersMutex.Unlock()

	conn.balanceMutex.Lock()
	conn.userAmounts[username] = 100
	conn.balanceMutex.Unlock()
}

func NewConnector() *Connector {
	conn := new(Connector)
	conn.users = make(map[string]string)
	conn.userAmounts = make(map[string]int)
	return conn
}
