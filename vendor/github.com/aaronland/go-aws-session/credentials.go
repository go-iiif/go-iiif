package session

import (
	"fmt"
	"strings"
)

func ValidCredentials() []string {

	valid := []string{
		"anon:",
		"env:",
		"iam:",
		"{PROFILE}",
		"{PATH}:{PROFILE}",
		"static:{ID}:{KEY}:{SECRET}",
	}

	return valid
}

func ValidCredentialsString() string {

	valid := ValidCredentials()
	return fmt.Sprintf("Valid credential flags are: %s", strings.Join(valid, ", "))
}
