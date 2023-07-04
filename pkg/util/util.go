package util

import (
	"log"
	"strconv"
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

func Float2string(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}
