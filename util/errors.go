package util

import (
	"errors"
	"os"
)

var logger = NewLog()

func FatalErrorCheck(err error) {
	if err != nil {
		logger.Printf(0, "\n%s\n", err.Error())
		os.Exit(1)
	}
}

func Assert(assertion bool, msg string) {
	if !assertion {
		FatalErrorCheck(errors.New("assertion failed: " + msg))
	}
}

func WarningErrorCheck(err error) bool {
	if err != nil {
		logger.Printf(0, "\n%s\n", err.Error())
		return false
	}

	return true
}
