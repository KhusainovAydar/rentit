package rentit

import "strconv"

func ParseUintOrDefault(str string) uint64 {
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

func ParseFloatOrDefault(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0.0
	}
	return val
}
