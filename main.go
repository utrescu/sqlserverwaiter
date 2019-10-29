package main

import (
	"errors"
	"fmt"
	"time"
)

func connectWithSQLServer() error {
	return errors.New("Connection failed")
}

// doItOrFail
func doItOrFail(timeout <-chan time.Time) error {

	tick := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-timeout:
			return errors.New("timed out")
		case <-tick:
			err := connectWithSQLServer()
			if err == nil {
				return nil
			}
			fmt.Println(err.Error())
		}
	}
}

func main() {
	//
	timeout := time.After(5 * time.Second)
	err := doItOrFail(timeout)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
	}

}
