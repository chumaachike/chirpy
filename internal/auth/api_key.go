package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authorization := headers.Get("authorization")
	if authorization == "" {
		return "", errors.New("authorization header missing")
	}
	const prefix = "ApiKey "
	if !strings.HasPrefix(authorization, prefix) {
		return "", errors.New("authorization header format must be Apikey <APIKEY>")
	}
	return strings.TrimSpace(authorization[len(prefix):]), nil
}
