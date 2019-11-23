package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"syscall"
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
			return errors.New("database ping fails")
		}
		return err
	}

	// When Collation is not null database is ready
	row := m.connection.QueryRow(fmt.Sprintf("SELECT DATABASEPROPERTYEX('%s', 'Collation') AS Collation", cmd.Database))
	var collation string
	err := row.Scan(&collation)
	if err != nil {
		return errors.New("database not ready")
	}

	// Database ready!!
	return nil
}

// NewConnection creates a connection with Sql Server with provided params
func NewConnection(host string, port int, database string, user string, password string, debug bool) (RepositoryReady, error) {

	query := url.Values{}
	query.Add("database", database)

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(user, password),
		Host:     fmt.Sprintf("%s:%d", host, port),
		RawQuery: query.Encode(),
	}

	if debug {
		fmt.Printf("DEBUG: %s\n", u.String())
	}

	db, err := sql.Open("sqlserver", u.String())
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

	connect, err := NewConnection(cmd.Server, cmd.Port, cmd.Database, cmd.User, cmd.Password, cmd.Debug)
	if err != nil {
		panic(fmt.Sprintf("Connection: %s", err.Error()))
	}

	timeout := time.After(cmd.Timeout)
	ok, err := doItOrFail(timeout, connect)
	if err != nil {
		fmt.Printf("Connection: %s\n", err.Error())
		syscall.Exit(1)
	} else {
		fmt.Printf("Connection: %v\n", ok)
	}

}
