package authentication

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/william22913/example-dating-apps/constanta"
	"github.com/william22913/example-dating-apps/custom_context"
	"github.com/william22913/example-dating-apps/custom_error"
	"github.com/william22913/example-dating-apps/token"
)

type authenticationUserAccess struct {
	jwtTokenValidator token.JWTTokenValidator
	redis             *redis.Client
}

func NewAuthenticationUserAccess(
	jwtTokenValidator token.JWTTokenValidator,
	redis *redis.Client,
) AuthenticationUserAccess {
	return authenticationUserAccess{
		jwtTokenValidator: jwtTokenValidator,
		redis:             redis,
	}
}

func (aua authenticationUserAccess) ValidateUserToken(
	ctx context.Context,
	header map[string]string,
) error {

	headerAuth := header["Authorization"]

	tokenModel, err := aua.jwtTokenValidator.ValidateToken(headerAuth)
	if err != nil {
		return err
	}

	redisResult, err := aua.redis.Get(headerAuth).Result()

	if err != nil {
		if err == redis.Nil {
			return custom_error.ErrUnauthorized
		}

		return err
	}

	var modelOnRedis token.RedisTokenModel

	err = json.Unmarshal([]byte(redisResult), &modelOnRedis)

	if err != nil {
		return custom_error.ErrUnauthorized
	}

	return aua.setUserToContext(ctx, tokenModel, modelOnRedis)
}

func (aua authenticationUserAccess) GenerateAndSaveUserToken(
	redismodel token.RedisTokenModel,
) (
	string,
	error,
) {

	token, err := aua.jwtTokenValidator.GenerateToken(redismodel.UserID, redismodel.UserUUID)
	if err != nil {
		return token, err
	}

	redisToken, err := json.Marshal(redismodel)
	if err != nil {
		return token, err
	}

	err = aua.redis.Set(token, redisToken, aua.jwtTokenValidator.GetTokenDuration()).Err()
	if err != nil {
		return token, err
	}

	return token, err
}

func (aua authenticationUserAccess) setUserToContext(
	ctx context.Context,
	tokenModel *token.PayloadJWTToken,
	modelOnRedis token.RedisTokenModel,
) error {

	_ctx, valid := ctx.Value(constanta.ApplicationContextConstanta).(*custom_context.ContextModel)
	if !valid {
		_ctx = custom_context.NewContextModel()
	}

	_ctx.AuthAccessTokenModel = *tokenModel
	_ctx.Account.UserID = modelOnRedis.UserID
	_ctx.Account.UserUUID = modelOnRedis.UserUUID
	_ctx.Account.FirstName = modelOnRedis.FirstName
	_ctx.Account.MiddleName = modelOnRedis.MiddleName
	_ctx.Account.LastName = modelOnRedis.LastName
	_ctx.Account.BirthDate, _ = time.Parse(modelOnRedis.BirthDate, modelOnRedis.BirthDate)
	_ctx.Account.PhoneNumber = modelOnRedis.PhoneNumber

	_ = context.WithValue(ctx, constanta.ApplicationContextConstanta, _ctx)

	return nil
}
