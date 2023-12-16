package api

import (
	"strings"
)

// extracts the bearer token from the Authorization header
func extractBearer(authorization string) string {
	var tokens = strings.Split(authorization, " ")
	if len(tokens) == 2 {
		return strings.Trim(tokens[1], " ")
	}
	return "0"
}

// return true if the user is not logged
func isNotLogged(token int) bool {
	if token == 0 {
		return true
	}

	return false
}
