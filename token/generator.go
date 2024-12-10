package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (j *jwtTokenValidator) GenerateToken(
	UserID int64,
	UUID string,
) (
	string,
	error,
) {
	payload := &PayloadJWTToken{
		UserID: UserID,
		StandardClaims: jwt.StandardClaims{
			Subject:   "user_token",
			Audience:  "user",
			Issuer:    "dating-apps",
			ExpiresAt: time.Now().Add(j.expired_time).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        UUID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return token.SignedString([]byte(j.secretKey))

}
