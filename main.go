package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/utrescu/sqlserverwait/cmd"
)

// RepositoryReady defines methods needed to this program
type RepositoryReady interface {
	IsAlive() error
}

// MsSQLConnection Defines a connection with Sql Server
type MsSQLConnection struct {
	connection *sql.DB
}

// IsAlive checks if Sql Server is up and accepts connections
func (m *MsSQLConnection) IsAlive() error {

	if err := m.connection.Ping(); err != nil {
		if err.Error() == "EOF" {
			return errors.New("database not ready")
		}
		return err
	}

	// When Collation is not null database is ready
	row := m.connection.QueryRow("SELECT DATABASEPROPERTYEX('BoIsBo', 'Collation') AS Collation")
	var collation string
	err := row.Scan(&collation)
	if err != nil {
		return errors.New("database not ready")
	}

	// Database ready!!
	return nil
}

// New creates a connection with Sql Server
func New(connectionString string) (RepositoryReady, error) {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, err
	}

	sqlconnection := &MsSQLConnection{
		connection: db,
	}

	return sqlconnection, nil
}

// doItOrFail tries until database is ready or time is over
func doItOrFail(timeout <-chan time.Time, connexio RepositoryReady) (bool, error) {

	tick := time.Tick(2 * time.Second)
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

	cmd.Execute()

	query := url.Values{}
	query.Add("database", cmd.Database)

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(cmd.User, cmd.Password),
		Host:     fmt.Sprintf("%s:%d", cmd.Server, cmd.Port),
		RawQuery: query.Encode(),
	}

	if cmd.Debug {
		fmt.Printf("DEBUG: %s\n", u.String())
	}

	connect, err := New(u.String())
	if err != nil {
		panic(fmt.Sprintf("Connection: %s", err.Error()))
	}

	timeout := time.After(cmd.Timeout)
	ok, err := doItOrFail(timeout, connect)
	if err != nil {
		fmt.Printf("Connection: %s\n", err.Error())
	} else {
		fmt.Printf("Connection: %v\n", ok)
	}

}
