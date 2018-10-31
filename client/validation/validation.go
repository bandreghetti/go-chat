package validation

import "regexp"

var userRegex = regexp.MustCompile("[:a-zA-Z:][:word:]+")

//ValidUsername checks given string against the userRegex definition
func ValidUsername(username string) bool {
	return userRegex.MatchString(username)
}
