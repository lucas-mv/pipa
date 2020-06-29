package utils

import (
	"encoding/json"
	"log"
	"strconv"
)

// FloatToString Converts floating point number to string
func FloatToString(num float64) string {
	return strconv.FormatFloat(num, 'f', 6, 64)
}

// PrettyPrint Prints idented json of interface i
func PrettyPrint(i interface{}) string {
	s, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	return string(s)
}
