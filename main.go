package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	server        = flag.String("server", "localhost", "Database server")
	port     *int = flag.Int("port", 1433, "Database port")
	user          = flag.String("user", "sa", "Database user")
	password      = flag.String("password", "X1nGuXunG1", "Database password")
	database      = flag.String("database", "BoIsBo", "Database name")
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

	flag.Parse()

	query := url.Values{}
	query.Add("database", *database)

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(*user, *password),
		Host:     fmt.Sprintf("%s:%d", *server, *port),
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
