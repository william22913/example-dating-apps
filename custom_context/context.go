package custom_context

import (
	"context"
	"time"

	"github.com/william22913/example-dating-apps/constanta"
	"github.com/william22913/example-dating-apps/token"
)

func (c *ContextModel) ToContext() context.Context {
	return context.WithValue(
		context.Background(),
		constanta.ApplicationContextConstanta,
		c,
	)
}

func NewContextModel() *ContextModel {
	ctx := new(ContextModel)
	return ctx
}

type ContextModel struct {
	ClientAccess         ClientAccess
	AuthAccessTokenModel token.PayloadJWTToken
	Account              Account
}

type ClientAccess struct {
	RequestID string
	Timestamp time.Time
	Headers   map[string]string
	Path      string
}

type Account struct {
	UserID      int64
	UserUUID    string
	FirstName   string
	MiddleName  string
	LastName    string
	BirthDate   time.Time
	PhoneNumber string
}
