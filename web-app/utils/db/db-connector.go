package db

import (
	"errors"
	//"github.com/lib/pq"
	"sync"
)

// this should be a database instead
// but for now everything is stored in memory
type Connector struct {
	Users                    map[string]string
	UserAmounts              map[string]int
	usersMutex, balanceMutex sync.Mutex
}

func (conn *Connector) GetPassword(username string) string {
	// what if user doesnt exist
	conn.usersMutex.Lock()
	defer conn.usersMutex.Unlock()
	return conn.Users[username]
}

func (conn *Connector) TransferTx(from, to string, amount int) error {
	// what if user doesnt exist
	conn.balanceMutex.Lock()
	defer conn.balanceMutex.Unlock()
	if conn.UserAmounts[from] < amount {
		return errors.New("Not enough funds on your account")
	}
	conn.UserAmounts[from] -= amount
	conn.UserAmounts[to] += amount
	return nil
}

func (conn *Connector) AddNewUser(username, password string) {
	conn.usersMutex.Lock()
	conn.Users[username] = password
	conn.usersMutex.Unlock()

	conn.balanceMutex.Lock()
	conn.UserAmounts[username] = 100
	conn.balanceMutex.Unlock()
}

func NewConnector() *Connector {
	conn := new(Connector)
	conn.Users = make(map[string]string)
	conn.UserAmounts = make(map[string]int)
	return conn
}
