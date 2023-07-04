package util

import "fmt"

func Logln(msg string) {
	fmt.Println(msg)
}

func Logf(format string, v ...any) {
	fmt.Printf(format, v...)
}
