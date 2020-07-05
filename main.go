package main

import (
	"fmt"
	"os"

	"github.com/utrescu/sqlserverwaiter/cmd"
	mssql "github.com/utrescu/sqlserverwaiter/mssql"
	ready "github.com/utrescu/sqlserverwaiter/ready"
)

func main() {

	cmd.Execute()

	connect, err := mssql.New(cmd.Server, cmd.Port, cmd.Database, cmd.User, cmd.Password, cmd.Debug)
	if err != nil {
		panic(fmt.Sprintf("Connection: %s", err.Error()))
	}

	if ready.Check(cmd.Timeout, connect) != nil {
		fmt.Printf("Connection: %s\n", err.Error())
		os.Exit(1)
	} else {
		fmt.Println("Ok")
	}
}
