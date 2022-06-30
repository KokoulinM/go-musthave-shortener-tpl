package helpers

import (
	"net/http"
)

func CreateCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}
}
