package util

import (
	"log"
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
