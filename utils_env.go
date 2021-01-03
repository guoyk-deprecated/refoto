package main

import (
	"os"
	"strconv"
	"strings"
)

func envStr(key string, output *string) (err error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value != "" {
		*output = value
	}
	return
}

func envBool(key string, output *bool) (err error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value != "" {
		var boolValue bool
		if boolValue, err = strconv.ParseBool(value); err != nil {
			return
		}
		*output = boolValue
	}
	return
}

func envInt(key string, output *int) (err error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value != "" {
		var intValue int
		if intValue, err = strconv.Atoi(value); err != nil {
			return
		}
		*output = intValue
	}
	return
}
