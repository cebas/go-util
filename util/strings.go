package util

import (
	"strconv"
	"strings"
)

func Float2string(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

func StringToFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f, err
}
