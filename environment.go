package gox

import "os"



func EnvReadStringOr(envIdentifier string, defaultValue string) string {
	value := os.Getenv(envIdentifier)

	if value == "" {
		return defaultValue
	}

	return value
}
