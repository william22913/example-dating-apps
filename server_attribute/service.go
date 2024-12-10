package server_attribute

import (
	loginService "github.com/william22913/example-dating-apps/service/login"
	signUpService "github.com/william22913/example-dating-apps/service/sign_up"
)

func (s *serverAttribute) InitService() {
	s.Services = Services{
		LoginService: loginService.NewLoginService(
			s.DAOs.UserDAO,
			s.PasswordGenerator,
			s.Auth,
		),
		SignUpService: signUpService.NewSignUpService(
			s.DAOs.UserDAO,
			s.PasswordGenerator,
		),
	}
}
