package misc

import (
	"strconv"
)

func Max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func Abs(a int32) int32 {
	if a < 0 {
		return -a
	}
	return a
}

func Atoi32(a string) int32 {
	v, err := strconv.Atoi(a)
	if err != nil {
		return 0
	}
	return int32(v)
}
