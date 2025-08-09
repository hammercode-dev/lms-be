package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

// VerifyToken verifies the given JWT token string using the configured secret key.
// It parses the token and returns the custom claims if the token is valid.
// Returns an error if the token is invalid or cannot be parsed.
//
// Parameters:
//   - token: the JWT token string to be verified
//
// Returns:
//   - *JwtCustomClaims: the custom claims extracted from the token if valid
//   - error: error if the token is invalid or parsing fails
func (j *jwtConfig) VerifyToken(token string) (*JwtCustomClaims, error) {
	tkn, err := jwt.ParseWithClaims(token, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return tkn.Claims.(*JwtCustomClaims), nil
}
