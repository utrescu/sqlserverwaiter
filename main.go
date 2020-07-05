package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/utrescu/sqlserverwaiter/cmd"
	ready "github.com/utrescu/sqlserverwaiter/ready"
	mssql "github.com/utrescu/sqlserverwaiter/mssql"
)


func main() {

	cmd.Execute()

	// Prepare SQL Connection
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

	connect, err := mssql.New(u.String(), cmd.Database)
	if err != nil {
		panic(fmt.Sprintf("Connection: %s", err.Error()))
	}

	result, err := ready.Check(cmd.Timeout, connect)
	if (err != nil) {
		fmt.Printf("Connection: %s\n", err.Error())
	}
	else
	{
		fmt.Println("Ok")
	}
}
