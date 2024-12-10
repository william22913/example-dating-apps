package authentication

import (
	"context"

	"github.com/william22913/example-dating-apps/token"
)

type AuthenticationUserAccess interface {
	ValidateUserToken(
		ctx context.Context,
		header map[string]string,
	) error

	GenerateAndSaveUserToken(
		redismodel token.RedisTokenModel,
	) (
		string,
		error,
	)
}
