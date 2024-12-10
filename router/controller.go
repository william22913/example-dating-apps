package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	// "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/william22913/example-dating-apps/config"
	http_validator "github.com/william22913/example-dating-apps/custom_endpoint"
	"github.com/william22913/example-dating-apps/util/endpoint/health"
)

var router *mux.Router

func InitHttpService(
	config config.Configuration,
	httpValidator *http_validator.HTTPController,
	health health.HealthEndpoint,
) {
	router = mux.NewRouter()
	httpValidator.Router(router)
	// router.HandleFunc("/metrics", promhttp.Handler().ServeHTTP).Methods(http.MethodGet)
	router.HandleFunc("/v1/health", health.CheckHealthConnection).Methods(http.MethodGet)

}

func StartService(
	config config.Configuration,
	httpValidator *http_validator.HTTPController,
) {

	log.Info().
		Str("action", "server.start").
		Int("port", config.Server.Port).
		Msg("HTTP Server Start.")

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port), router)

	if err != nil {
		log.Fatal().
			Str("action", "server.stop").
			Int("port", config.Server.Port).
			Err(err).
			Msg("HTTP Server Stopped.")
	}
}
