// Package helpers contains auxiliary functions
package helpers

import (
	"net/http"
)

// CreateCookie a represents an HTTP cookie as sent in the Set-Cookie header of an
func CreateCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}
}
