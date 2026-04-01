package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is required")
	}

	const bearerScheme = "Bearer "
	if !strings.HasPrefix(authHeader, bearerScheme) {
		return "", fmt.Errorf("invalid authorization header format")
	}

	tokenString := strings.TrimPrefix(authHeader, bearerScheme)
	if tokenString == "" {
		return "", fmt.Errorf("token is required")
	}

	return tokenString, nil
}
