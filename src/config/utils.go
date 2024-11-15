package config

import (
	"fmt"
	"os"
	"strconv"
)

// gets a string value from the environment
// convenience function, guarantees non-empty value or error
func fromEnvString(k string) (string, error) {
	v := os.Getenv(k)
	if v == "" {
		return "", fmt.Errorf("environment variable %s was empty", k)
	}

	return v, nil
}

// gets a number out of an env variable.
// convenience function, guarantees non-empty value or error
func fromEnvInt(k string) (int, error) {
	v, err := fromEnvString(k)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("unable to parse %s as int: %+v", v, err)
	}

	return i, nil
}
