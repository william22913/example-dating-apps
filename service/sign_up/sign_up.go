package sign_up

import (
	"database/sql"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/william22913/example-dating-apps/custom_context"
	"github.com/william22913/example-dating-apps/custom_error"
	"github.com/william22913/example-dating-apps/dao/user"
	"github.com/william22913/example-dating-apps/dto/in"
	"github.com/william22913/example-dating-apps/password"
	"github.com/william22913/example-dating-apps/repository"
)

func NewSignUpService(
	userDAO user.UserDAO,
	passwordGenerator password.PasswordAlgorithm,

) SignUpService {
	regexPhone := regexp.MustCompile(`[+][0-9]+[-][1-9][0-9]{8,12}$`)
	regexPassword := regexp.MustCompile(`^[A-Za-z0-9]+$`)

	return &signUpService{
		userDAO:           userDAO,
		regexPhone:        regexPhone,
		regexPassword:     regexPassword,
		passwordGenerator: passwordGenerator,
	}
}

type signUpService struct {
	userDAO           user.UserDAO
	regexPhone        *regexp.Regexp
	regexPassword     *regexp.Regexp
	passwordGenerator password.PasswordAlgorithm
}

func (s *signUpService) SignUpUser(
	ctx *custom_context.ContextModel,
	dto interface{},
) (
	header map[string]string,
	result interface{},
	err error,
) {
	dtoIn := dto.(*in.SignUpDTOIn)

	err = s.validateIncomingRequest(dtoIn)
	if err != nil {
		return
	}

	salt := GetUUID()
	password, err := s.passwordGenerator.HidePassword(dtoIn.Password, salt)
	if err != nil {
		return
	}

	dataOnDB, err := s.userDAO.CheckIsPhoneExist(repository.User{
		PhoneNumber: sql.NullString{String: dtoIn.PhoneNumber},
	})

	if err != nil || dataOnDB.ID.Valid {
		err = custom_error.ErrDataUsed.Param("PHONE_NUMBER")
		return
	}

	err = s.userDAO.InsertUserData(
		repository.CompletedUserData{
			User: repository.User{
				PhoneNumber: sql.NullString{String: dtoIn.PhoneNumber},
				Password:    sql.NullString{String: password},
				Gender:      sql.NullString{String: dtoIn.Gender},
				FirstName:   sql.NullString{String: dtoIn.FirstName},
				MiddleName:  sql.NullString{String: dtoIn.MiddleName},
				LastName:    sql.NullString{String: dtoIn.LastName},
				BirthDate:   sql.NullTime{Time: dtoIn.BirthDate},
			},
			UserPreferences: repository.UserPreferences{
				Gender: sql.NullString{String: dtoIn.Gender},
				MinAge: sql.NullInt64{Int64: int64(dtoIn.Preferences.MinAge)},
				MaxAge: sql.NullInt64{Int64: int64(dtoIn.Preferences.MaxAge)},
			},
			UserPassions: repository.UserPassions{
				Tags: dtoIn.Passions,
			},
			Salt: repository.Salt{
				SaltKey: sql.NullString{String: salt},
			},
		},
	)

	return
}

func (s *signUpService) GetDTO() interface{} {
	return &in.SignUpDTOIn{}
}

func (s *signUpService) validateIncomingRequest(dtoIn *in.SignUpDTOIn) error {

	if !s.regexPassword.MatchString(dtoIn.Password) {
		return custom_error.ErrValidationBody.Param("Password", "Regex")
	}

	if !s.regexPhone.MatchString(dtoIn.PhoneNumber) {
		return custom_error.ErrValidationBody.Param("PhoneNumber", "Regex")
	}

	dtoIn.BirthDate, _ = time.Parse("2006-01-02", dtoIn.BirthDateStr)

	return nil
}

func GetUUID() (output string) {
	UUID, _ := uuid.NewRandom()
	output = UUID.String()
	output = strings.Replace(output, "-", "", -1)
	return
}
