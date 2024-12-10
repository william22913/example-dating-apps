package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type PayloadJWTToken struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

type JWTTokenValidator interface {
	ValidateToken(
		token string,
	) (
		*PayloadJWTToken,
		error,
	)

	GenerateToken(
		UserID int64,
		UUID string,
	) (
		token string,
		err error,
	)

	GetTokenDuration() time.Duration
}

type RedisTokenModel struct {
	UserID      int64  `json:"uid"`
	UserUUID    string `json:"uuid"`
	FirstName   string `json:"fn"`
	MiddleName  string `json:"mn"`
	LastName    string `json:"ln"`
	BirthDate   string `json:"bd"`
	PhoneNumber string `json:"pn"`
	Gender      string `json:"gdr"`
}
