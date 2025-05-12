package jwt

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hammer-code/lms-be/domain"
)

type (
	jwtConfig struct {
		SecretKey string
	}

	JwtCustomClaims struct {
		ID       int    `json:"id"`
		UserName string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		jwt.RegisteredClaims
	}

	JWT interface {
		GenerateAccessToken(c context.Context, user *domain.User, expiredTime int) (*string, error)
		VerifyToken(token string) (*JwtCustomClaims, error)
	}
)

func NewJwt(secretKey string) JWT {
	return &jwtConfig{
		SecretKey: secretKey,
	}
}
