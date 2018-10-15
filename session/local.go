package session

import (
	"fmt"
	"github.com/golang/glog"
	"time"
)

const (
	TimeoutSec = 5
)

func LocalSession(command <-chan string, output chan<- string, disconnect <-chan bool) {
	select {
	case c := <-command:
		time.Sleep(1 * time.Second) // test sleep
		output <- fmt.Sprintf("Echo command:\n  %s", c)

	case d := <-disconnect:
		if d {
			glog.Info("Received disconnect request. Return")
			close(output)

			return
		}

	case <-time.After(TimeoutSec * time.Second):
		glog.Infof("Timeout was reach. Drop session")
		close(output)

		return
	}
}
