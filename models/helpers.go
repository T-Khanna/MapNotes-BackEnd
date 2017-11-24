package models

import (
	"strconv"
	"strings"
)

func ConvertIntArrayToString(array []int64) string {
	var valuesText = []string{}
	// Create a string slice using strconv.Itoa.
	// ... Append string
	for i := range array {
		number := array[i]
		text := strconv.FormatInt(number, 10)
		valuesText = append(valuesText, text)
	}
	// Join our string slice.
	result := "(" + strings.Join(valuesText, ",") + ")"
	return result
}
