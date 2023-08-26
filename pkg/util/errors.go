package util

import (
	"fmt"
	"log"
	"time"
)

func FatalErrorCheck(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func WarningErrorCheck(err error) bool {
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

// Retry will loop upto 'attempts times with a logaritmical backoff
func Retry(attempts int, sleepTime time.Duration, f func() error) (err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			Logf("retrying after error: %v", err)
			time.Sleep(sleepTime)
			sleepTime *= 2
		}
		err = f()
		if err == nil {
			return
		}
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
