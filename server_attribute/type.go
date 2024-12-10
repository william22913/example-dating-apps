package server_attribute

import (
	"database/sql"

	"github.com/go-redis/redis/v7"
	"github.com/william22913/example-dating-apps/authentication"
	"github.com/william22913/example-dating-apps/bundles"
	"github.com/william22913/example-dating-apps/config"
	"github.com/william22913/example-dating-apps/custom_endpoint"
	"github.com/william22913/example-dating-apps/dao/user"
	"github.com/william22913/example-dating-apps/password"
	loginService "github.com/william22913/example-dating-apps/service/login"
	signUpService "github.com/william22913/example-dating-apps/service/sign_up"
	"github.com/william22913/example-dating-apps/token"
	"github.com/william22913/example-dating-apps/util/endpoint/health"
)

type serverAttribute struct {
	config            config.Configuration
	Postgresql        *sql.DB
	RedisClient       *redis.Client
	Auth              authentication.AuthenticationUserAccess
	TokenValidator    token.JWTTokenValidator
	Bundles           bundles.Bundles
	HttpController    *custom_endpoint.HTTPController
	DAOs              DAOs
	Services          Services
	PasswordGenerator password.PasswordAlgorithm
	HealthCheck       health.HealthEndpoint
}

type Services struct {
	LoginService  loginService.LoginService
	SignUpService signUpService.SignUpService
}

type DAOs struct {
	UserDAO user.UserDAO
}
