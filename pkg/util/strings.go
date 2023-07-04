package util

import (
	"strconv"
)

func Float2string(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}
