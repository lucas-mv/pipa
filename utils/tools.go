package utils

import (
	"strconv"
)

// FloatToString Converts floating point number to string
func FloatToString(num float64) string {
	return strconv.FormatFloat(num, 'f', 6, 64)
}
