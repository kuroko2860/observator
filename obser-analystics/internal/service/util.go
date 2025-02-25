package service

import (
	"fmt"
	"strconv"
)

func ParseFromToStringToInt(from, to string) (int64, int64) {
	fromInt, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		fmt.Println("from err")
	}
	toInt, err := strconv.ParseInt(to, 10, 64)
	if err != nil {
		fmt.Println("to err")
	}
	return fromInt, toInt
}
func ParseUnitToInterval(unit string) int64 {
	switch unit {
	case "second":
		return 1000
	case "minute":
		return 60 * 1000
	case "hour":
		return 60 * 60 * 1000
	case "day":
		return 24 * 60 * 60 * 1000
	default:
		return 60 * 60 * 1000

	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
