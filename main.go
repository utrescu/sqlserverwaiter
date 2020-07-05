package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/utrescu/sqlserverwait/cmd"
	database "github.com/utrescu/sqlserverwait/db"
	mssql "github.com/utrescu/sqlserverwait/mssql"
)

// doItOrFail tries until database is ready or time is over
func doItOrFail(timeout <-chan time.Time, connexio database.RepositoryReady) (bool, error) {

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

	connect, err := mssql.New(u.String())
	if err != nil {
		panic(fmt.Sprintf("Connection: %s", err.Error()))
	}

	timeout := time.After(cmd.Timeout)
	ok, err := doItOrFail(timeout, connect)
	if err != nil {
		fmt.Printf("Connection: %s\n", err.Error())
		os.Exit(1)
	} else {
		fmt.Printf("Connection: %v\n", ok)
	}

}
