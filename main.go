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

func connectWithSQLServer() error {

	query := url.Values{}
	query.Add("database", database)

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(user, password),
		Host:     fmt.Sprintf("%s:%d", server, port),
		RawQuery: query.Encode(),
	}
	fmt.Println(u.String())
	// connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", server, user, password, port, database)
	db, err := sql.Open("sqlserver", u.String())
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	row := db.QueryRow("SELECT DATABASEPROPERTYEX('BoIsBo', 'Collation') AS Collation")
	var collation string
	err = row.Scan(&collation)
	if err != nil {
		return err
	}
	fmt.Println(collation)
	return nil
}

// doItOrFail
func doItOrFail(timeout <-chan time.Time) (bool, error) {

	tick := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-timeout:
			return false, errors.New("timed out")
		case <-tick:
			err := connectWithSQLServer()
			if err == nil {
				return true, nil
			}
			fmt.Println(err.Error())
		}
	}
}

func main() {
	//
	timeout := time.After(30 * time.Second)
	ok, err := doItOrFail(timeout)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
	}

	fmt.Printf("Connection: %v", ok)

}
