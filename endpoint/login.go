package endpoint

import (
	"net/http"

	"github.com/william22913/example-dating-apps/authentication"
	"github.com/william22913/example-dating-apps/custom_endpoint"
	"github.com/william22913/example-dating-apps/service/login"
)

func NewLoginEndpoint(
	httpController *custom_endpoint.HTTPController,
	loginService login.LoginService,
	auth authentication.AuthenticationUserAccess,
) *loginEndpoint {
	return &loginEndpoint{
		httpController: httpController,
		loginService:   loginService,
		auth:           auth,
	}
}

type loginEndpoint struct {
	httpController *custom_endpoint.HTTPController
	loginService   login.LoginService
	auth           authentication.AuthenticationUserAccess
}

func (e *loginEndpoint) RegisterEndpoint() {
	e.httpController.HandleFunc(
		custom_endpoint.NewHandleFuncParam(
			"/v1/appdate/login",
			e.httpController.WrapService(
				custom_endpoint.NewWarpServiceParam(
					e.loginService,
					e.loginService.Login,
					nil,
				),
			),
			http.MethodPost,
		),
	)
}
