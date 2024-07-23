package util

import (
	"errors"
	"strconv"
)

// parse number in string to int64
func ParseStringToInt64(s string) (int64, error) {
	if s == "" {
		return 0, errors.New("empty string")
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// parse number in string to int32
func ParseStringToInt32(s string) (int32, error) {
	if s == "" {
		return 0, errors.New("empty string")
	}
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}
