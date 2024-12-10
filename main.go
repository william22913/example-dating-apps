package main

import (
	"github.com/rs/zerolog/log"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/william22913/example-dating-apps/config"
	"github.com/william22913/example-dating-apps/router"
	"github.com/william22913/example-dating-apps/server_attribute"
	dbmigration "github.com/william22913/example-dating-apps/util/db_migration"
)

func main() {
	cfg := config.AppConfig
	serverAttribute := server_attribute.NewServerAttribute(cfg)
	err := serverAttribute.Init()

	if err != nil {
		log.Fatal().Err(err).Msg("error found when init server attribute")
	}

	_, err = dbmigration.NewDBMigration().DbMigratePostresql(serverAttribute.Postgresql, "sql_migrations", cfg.Postgresql.DefaultSchema, migrate.Up)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	router.InitHttpService(
		cfg,
		serverAttribute.HttpController,
		serverAttribute.HealthCheck,
	)

	serverAttribute.InitEndpoint()

	router.StartService(cfg, serverAttribute.HttpController)
}
