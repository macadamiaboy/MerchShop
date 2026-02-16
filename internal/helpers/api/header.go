package api

import (
	"fmt"
	"net/http"
)

func GetToken(r *http.Request) (string, error) {
	env := "helpers.api.header.GetToken"

	token := r.Header.Get("Authorization")

	if token == "" {
		return "", fmt.Errorf("%s: the token is not provided", env)
	} else {
		return token, nil
	}
}
