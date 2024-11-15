package config

import (
	"fmt"
	"os"
	"strconv"
)

func fromEnvString(k string) (string, error) {
	v := os.Getenv(k)
	if v == "" {
		return "", fmt.Errorf("environment variable %s was empty", k)
	}

	return v, nil
}

func fromEnvUint(k string) (uint, error) {
	v, err := fromEnvString(k)
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse %s as int: %+v", v, err)
	}

	return uint(i), nil
}
