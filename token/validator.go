package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/william22913/example-dating-apps/custom_error"
)

func NewJWTTokenValidator(
	secretKey string,
	expired_time time.Duration,
) JWTTokenValidator {
	return &jwtTokenValidator{
		secretKey:    secretKey,
		expired_time: expired_time,
	}
}

type jwtTokenValidator struct {
	secretKey    string
	expired_time time.Duration
}

func (j *jwtTokenValidator) ValidateToken(token string) (*PayloadJWTToken, error) {
	// Parse the token
	parsedToken, err := jwt.ParseWithClaims(token, &PayloadJWTToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	payload, ok := parsedToken.Claims.(*PayloadJWTToken)
	if !ok || !parsedToken.Valid {
		return nil, custom_error.ErrUnauthorized
	}

	// Check if the token is expired
	if payload.ExpiresAt < time.Now().Unix() {
		return nil, custom_error.ErrExpiredToken
	}

	return payload, nil
}

func (j *jwtTokenValidator) GetTokenDuration() time.Duration {
	return j.expired_time
}
