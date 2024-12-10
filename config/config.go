package config

import (
	"encoding/json"
	"fmt"
	logger "log"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/kelseyhightower/envconfig"
)

var (
	AppConfig Configuration
)

type Configuration struct {
	Log        logs       `envconfig:"log"`
	Postgresql Postgresql `envconfig:"postgresql"`
	Server     server     `envconfig:"server"`
	Token      token      `envconfig:"token"`
	Redis      redis      `envconfig:"redis"`
}

type logs struct {
	Level int8 `envconfig:"level" default:"0"`
}

type token struct {
	UserKey  string        `envconfig:"user_key" default:"user_key"`
	Duration time.Duration `envconfig:"duration" default:"24h"`
}

type server struct {
	Version string `envconfig:"version" default:"1.0.0"`
	Host    string `envconfig:"host"`
	Port    int    `envconfig:"port" default:"8843"`
}

type Postgresql struct {
	Host              string `envconfig:"host" default:"localhost"`
	Port              int    `envconfig:"port" default:"5433"`
	Username          string `envconfig:"username" default:"postgres"`
	Password          string `envconfig:"password" default:"test"`
	DBName            string `envconfig:"dbname" default:"dating_apps"`
	SSLMode           string `envconfig:"sslmode" default:"false"`
	DefaultSchema     string `envconfig:"defaultschema" default:"dating_apps"`
	MaxOpenConnection int    `envconfig:"maxopenconnection" default:"500"`
	MaxIdleConnection int    `envconfig:"maxidleconnection" default:"500"`
}

type redis struct {
	Host       string `envconfig:"host" default:"localhost"`
	Port       int    `envconfig:"port" default:"6379"`
	DB         int    `envconfig:"db" default:"0"`
	Password   string `envconfig:"password"`
	Username   string `envconfig:"username"`
	MaxRetries int    `envconfig:"max_retries"`
}

func init() {

	if err := envconfig.Process(
		"dating_apps",
		&AppConfig,
	); err != nil {
		logger.Fatal(err)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Logger = logger

	PrintConfig(AppConfig)
}

func PrintConfig(c Configuration) {
	data, _ := json.MarshalIndent(c, "", "\t")
	fmt.Println(string(data))
}
