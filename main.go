package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	server   = "localhost"
	user     = "sa"
	port     = 1433
	password = "X1nGuXunG1"
	database = "BoIsBo"
)

// Repository defineix els mètodes de comprovació que hi ha d'haver
type Repository interface {
	IsAlive() error
}

// ServerConnection Defineix una connexió amb SQLServer
type ServerConnection struct {
	connection *sql.DB
}

// IsAlive comprova si hi ha connexió amb la base de dades i si està activa
func (m *ServerConnection) IsAlive() error {

	if err := m.connection.Ping(); err != nil {
		if err.Error() == "EOF" {
			return errors.New("Database not ready")
		}
		return errors.New("Unable to open connection")
	}

	// When Collation is not null database is ready
	row := m.connection.QueryRow("SELECT DATABASEPROPERTYEX('BoIsBo', 'Collation') AS Collation")
	var collation string
	err := row.Scan(&collation)
	if err != nil {
		return errors.New("Database not ready")
	}

	// Database ready!!
	return nil
}

// New crea una connexió amb la base de dades
func New(connectionString string) (Repository, error) {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, err
	}

	sqlconnection := &ServerConnection{
		connection: db,
	}

	return sqlconnection, nil
}

// doItOrFail
func doItOrFail(timeout <-chan time.Time, connexio Repository) (bool, error) {

	tick := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-timeout:
			return false, errors.New("timed out")
		case <-tick:
			err := connexio.IsAlive()
			if err == nil {
				return true, nil
			}
			fmt.Printf(".. %s\n", err.Error())
		}
	}
}

func main() {

	query := url.Values{}
	query.Add("database", database)

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(user, password),
		Host:     fmt.Sprintf("%s:%d", server, port),
		RawQuery: query.Encode(),
	}

	database, err := New(u.String())
	if err != nil {
		panic(fmt.Sprintf("Connection: %s", err.Error()))
	}

	timeout := time.After(30 * time.Second)
	ok, err := doItOrFail(timeout, database)
	if err != nil {
		fmt.Printf("Connection: %s\n", err.Error())
	} else {
		fmt.Printf("Connection: %v\n", ok)
	}

}
