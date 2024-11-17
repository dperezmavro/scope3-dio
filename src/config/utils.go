package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// gets a string value from the environment
// convenience function, guarantees non-empty value or error
func fromEnvStringDefault(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}

	return v
}

// gets a string value from the environment
// convenience function, guarantees non-empty value or error
func fromEnvBoolDefault(k string, d bool) bool {
	v := os.Getenv(k)
	if strings.ToLower(v) == "true" {
		return true
	}

	return d
}

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
func fromEnvIntDefault(k string, d int) (int, error) {
	v, err := fromEnvString(k)
	if err != nil {
		return d, err
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("unable to parse %s as int: %+v", v, err)
	}

	return i, nil
}
