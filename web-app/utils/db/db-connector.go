package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// this should be a database instead
// but for now everything is stored in memory
type Connector struct {
	pgConn *pgxpool.Pool
}

func (conn *Connector) GetPassword(username string) string {
	var password string
	err := conn.pgConn.QueryRow(context.Background(), "select password from users where username=$1", username).Scan(&password)
	if err != nil {
		return ""
	}

	return password
}

func (conn *Connector) AddNewUser(username, password string) error {
	_, err := conn.pgConn.Exec(context.Background(), "insert into users(id, username, password) values(gen_random_uuid(), $1, $2)", username, password)
	if err != nil {
		return err
	}

	return nil
}

func NewConnector() (*Connector, error) {
	conn := new(Connector)
	url := "postgres://postgres:password123@db:5432/postgres"
	pgConn, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}

	_, err = pgConn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS users (id uuid NOT NULL, username text UNIQUE NOT NULL, password text NOT NULL);")
	if err != nil {
		return nil, err
	}

	conn.pgConn = pgConn

	return conn, nil
}

func WaitForDb() (*Connector, error) {
	for trial := 0; trial < 3; trial++ {
		connector, err := NewConnector()
		if err == nil {
			return connector, nil
		}
		time.Sleep(time.Second * 10)
	}

	return nil, fmt.Errorf("couldn't connect to the database")
}
