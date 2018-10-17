package utilx

import "os"



func EnvReadStringOr(envIdentifier string, defaultValue string) string {
	value := os.Getenv(envIdentifier)

	if value == "" {
		return defaultValue
	}

	return value
}



func EnvReadBoolOr(envIdentifier string, defaultValue bool) bool {
	value := os.Getenv(envIdentifier)

	if value == "" {
		return defaultValue
	}

	if value == "false" {
		return false
	}


	if value == "true" {
		return true
	}

	return defaultValue
}

