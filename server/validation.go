package main

import "regexp"

var userRegex = regexp.MustCompile("[:a-zA-Z:][:word:]+")

func validUsername(username string) bool {
	if username == "server" {
		return false
	}
	return userRegex.MatchString(username)
}
