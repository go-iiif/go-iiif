package session

import (
	"fmt"
	"strings"
)

const AnonymousCredentialsString string = "anon:"

const EnvironmentCredentialsString string = "env:"

const IAMCredentialsString string = "iam:"

const ProfileCredentialsTemplate string = "{profile}"

const ProfileFromPathCredentialsTemplate string = "{path}:{profile}"

const StaticCredentialsTemplate string = "static:{id}:{key}:{secret}"

const StaticCredentialsPrefix string = "static:"

func ValidCredentials() []string {

	valid := []string{
		AnonymousCredentialsString,
		EnvironmentCredentialsString,
		IAMCredentialsString,
		ProfileCredentialsTemplate,
		ProfileFromPathCredentialsTemplate,
		StaticCredentialsTemplate,
	}

	return valid
}

func ValidCredentialsString() string {

	valid := ValidCredentials()
	return fmt.Sprintf("Valid credential flags are: %s", strings.Join(valid, ", "))
}
