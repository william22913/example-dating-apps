package server_attribute

import (
	"github.com/rs/zerolog/log"

	"github.com/go-playground/validator/v10"
	"github.com/william22913/example-dating-apps/authentication"
	"github.com/william22913/example-dating-apps/bundles"
	"github.com/william22913/example-dating-apps/config"
	"github.com/william22913/example-dating-apps/custom_endpoint"
	"github.com/william22913/example-dating-apps/custom_error"
	"github.com/william22913/example-dating-apps/dao/user"
	"github.com/william22913/example-dating-apps/password"
	"github.com/william22913/example-dating-apps/token"
	"github.com/william22913/example-dating-apps/util"
	"github.com/william22913/example-dating-apps/util/endpoint/health"
)

func NewServerAttribute(
	config config.Configuration,
) serverAttribute {
	return serverAttribute{
		config: config,
	}
}

func (s *serverAttribute) Init() (
	err error,
) {

	s.Postgresql = util.GetDbConnection(
		util.DBAddressParam().
			Host(s.config.Postgresql.Host).
			Port(s.config.Postgresql.Port).
			DBName(s.config.Postgresql.DBName).
			DefaultSchema(s.config.Postgresql.DefaultSchema).
			Username(s.config.Postgresql.Username).
			Password(s.config.Postgresql.Password),
	)

	s.RedisClient = util.ConnectRedis(
		util.NewRedisParam(
			s.config.Redis.Host,
			s.config.Redis.Port,
		).
			DB(s.config.Redis.DB).
			MaxRetries(s.config.Redis.MaxRetries).
			Password(s.config.Redis.Password).
			Username(s.config.Redis.Username),
	)

	s.Bundles, err = bundles.NewBundles("i18n", "en-US")
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to create bundles")
		return
	}

	errFormator := custom_error.NewErrorFormator(
		s.Bundles,
	).DefaultInternalCode("E-5-DATEAPPS-SRV-001").
		DefaultLanguage("en-US").
		Version(s.config.Server.Version)

	s.HttpController = custom_endpoint.NewHTTPController(validator.New())
	s.HttpController.Formator(errFormator)
	s.HttpController.Version(s.config.Server.Version)
	s.TokenValidator = token.NewJWTTokenValidator(s.config.Token.UserKey, s.config.Token.Duration)
	s.PasswordGenerator = password.NewBcryptPasswordAlgorithm()
	s.Auth = authentication.NewAuthenticationUserAccess(s.TokenValidator, s.RedisClient)

	s.DAOs = DAOs{
		UserDAO: user.NewPostgresqlUserDAO(s.Postgresql),
	}

	s.HealthCheck = health.NewHealthEndpoint(
		health.NewListTools("db", health.NewDBConnChecker(s.Postgresql)),
		health.NewListTools("redis", health.NewRedisConnChecker(s.RedisClient.Conn())),
	)

	s.InitService()

	return nil
}
