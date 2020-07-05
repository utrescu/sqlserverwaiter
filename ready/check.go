package ready

import (
	"errors"
	"fmt"
	"time"
)

// RepositoryReady defines methods needed to this program
type RepositoryReady interface {
	IsAlive() error
}

func doItOrFail(timeout <-chan time.Time, connexio RepositoryReady) error {

	tick := time.Tick(2 * time.Second)
	for {
		select {
		case <-timeout:
			return errors.New("timed out")
		case <-tick:
			err := connexio.IsAlive()
			if err == nil {
				return nil
			}
			fmt.Printf(".. %s\n", err.Error())
		}
	}
}

// Check determines when the repository is ready or timeouts
func Check(maxTime time.Duration, connection RepositoryReady) error {
	timeout := time.After(maxTime)
	err := doItOrFail(timeout, connection)
	if err != nil {
		return err
	}
	return nil
}
