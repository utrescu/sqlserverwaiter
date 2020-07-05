package ready

import (
	"errors"
	"fmt"
	"time"
)

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

// Check determines when the repository is ready or timeouts
func Check(maxTime time.Duration, connection RepositoryReady) (bool, error) {
	timeout := time.After(maxTime)
	ok, err := doItOrFail(timeout, connection)
	if err != nil {
		return false, err
	}
	return true, nil
}
