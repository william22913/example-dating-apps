package endpoint

import (
	"net/http"

	"github.com/william22913/example-dating-apps/custom_endpoint"
	"github.com/william22913/example-dating-apps/service/sign_up"
)

func NewSignUpEndpoint(
	httpController *custom_endpoint.HTTPController,
	signUpService sign_up.SignUpService,
) *signUpEndpoint {
	return &signUpEndpoint{
		httpController: httpController,
		signUpService:  signUpService,
	}
}

type signUpEndpoint struct {
	httpController *custom_endpoint.HTTPController
	signUpService  sign_up.SignUpService
}

func (e *signUpEndpoint) RegisterEndpoint() {
	e.httpController.HandleFunc(
		custom_endpoint.NewHandleFuncParam(
			"/v1/appdate/register",
			e.httpController.WrapService(
				custom_endpoint.NewWarpServiceParam(
					e.signUpService,
					e.signUpService.SignUpUser,
					nil,
				),
			),
			http.MethodPost,
		),
	)
}
