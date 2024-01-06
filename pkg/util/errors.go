package util

import (
	"os"
)

func FatalErrorCheck(err error) {
	if err != nil {
		Logln(err.Error())
		os.Exit(1)
	}
}

func WarningErrorCheck(err error) bool {
	if err != nil {
		Logln(err.Error())
		return false
	}

	return true
}
