package login

import (
	"database/sql"
	"regexp"

	"github.com/william22913/example-dating-apps/authentication"
	"github.com/william22913/example-dating-apps/custom_context"
	"github.com/william22913/example-dating-apps/custom_error"
	"github.com/william22913/example-dating-apps/dao/user"
	"github.com/william22913/example-dating-apps/dto/in"
	"github.com/william22913/example-dating-apps/password"
	"github.com/william22913/example-dating-apps/repository"
	"github.com/william22913/example-dating-apps/token"
)

func NewLoginService(
	userDAO user.UserDAO,
	passwordGenerator password.PasswordAlgorithm,
	auth authentication.AuthenticationUserAccess,
) LoginService {
	regexPhone := regexp.MustCompile(`[+][0-9]+[-][1-9][0-9]{8,12}$`)

	return &loginService{
		userDAO:           userDAO,
		regexPhone:        regexPhone,
		passwordGenerator: passwordGenerator,
		auth:              auth,
	}
}

type loginService struct {
	userDAO           user.UserDAO
	regexPhone        *regexp.Regexp
	passwordGenerator password.PasswordAlgorithm
	auth              authentication.AuthenticationUserAccess
}

func (s *loginService) Login(
	ctx *custom_context.ContextModel,
	dto interface{},
) (
	header map[string]string,
	result interface{},
	err error,
) {
	dtoIn := dto.(*in.LoginDTOIn)

	err = s.validateIncomingRequest(dtoIn)
	if err != nil {
		return
	}

	userData, err := s.userDAO.GetUserDataForLogin(
		repository.User{
			PhoneNumber: sql.NullString{String: dtoIn.PhoneNumber},
		},
	)

	if err != nil {
		return
	}

	if !userData.User.ID.Valid {
		err = custom_error.ErrUnauthorized
		return
	}

	if !s.passwordGenerator.CheckPassword(
		dtoIn.Password,
		userData.Salt.SaltKey.String,
		userData.User.Password.String,
	) {
		err = custom_error.ErrUnauthorized
		return
	}

	token, err := s.auth.GenerateAndSaveUserToken(token.RedisTokenModel{
		UserID:      userData.User.ID.Int64,
		UserUUID:    userData.User.UUIDKey.String,
		FirstName:   userData.User.FirstName.String,
		MiddleName:  userData.User.MiddleName.String,
		LastName:    userData.User.LastName.String,
		BirthDate:   userData.User.BirthDate.Time.Format("2006-01-02"),
		PhoneNumber: dtoIn.PhoneNumber,
	})

	header = map[string]string{
		"Authorization": token,
	}

	result = token

	return
}

func (s *loginService) GetDTO() interface{} {
	return &in.LoginDTOIn{}
}

func (s *loginService) validateIncomingRequest(dtoIn *in.LoginDTOIn) error {

	if !s.regexPhone.MatchString(dtoIn.PhoneNumber) {
		return custom_error.ErrValidationBody.Param("PhoneNumber", "Regex")
	}

	return nil
}
