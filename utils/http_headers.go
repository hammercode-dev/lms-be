package utils

import (
	"net/http"
	"strings"
)

// ExtractBearerToken extracts the Bearer token from the Authorization header of an HTTP request.
// If the Authorization header does not start with "Bearer ", it returns the header value as is.
// If the header starts with "Bearer ", it returns the token part after "Bearer ".
//
// Parameters:
//   - r: pointer to http.Request from which to extract the token
//
// Returns:
//   - *string: pointer to the extracted token string
func ExtractBearerToken(r *http.Request) *string {
	token := r.Header.Get("Authorization")
	if !strings.HasPrefix(token, "Bearer ") {
		return &token
	}

	token = strings.Split(token, "Bearer ")[1]
	return &token
}
