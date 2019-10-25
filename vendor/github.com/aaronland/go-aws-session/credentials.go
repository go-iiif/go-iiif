package session

import (
	"fmt"
	"strings"
)

func ValidCredentials() []string {

	valid := []string{
		"env:",
		"iam:",
		"{PROFILE}",
		"{PATH}:{PROFILE}",
	}

	return valid
}

func ValidCredentialsString() string {

	valid := ValidCredentials()
	return fmt.Sprintf("Valid credential flags are: %s", strings.Join(valid, ", "))
}
