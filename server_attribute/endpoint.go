package server_attribute

import (
	endpoints "github.com/william22913/example-dating-apps/endpoint"
	endpointUtil "github.com/william22913/example-dating-apps/util/endpoint"
)

func (s *serverAttribute) InitEndpoint() {
	endpoint := endpointUtil.NewEndpoint()

	endpoint.AddEndpoint(
		endpoints.NewLoginEndpoint(
			s.HttpController,
			s.Services.LoginService,
			s.Auth,
		),
	)

	endpoint.AddEndpoint(
		endpoints.NewSignUpEndpoint(
			s.HttpController,
			s.Services.SignUpService,
		),
	)

	endpoint.ServeEndpoint()
}
