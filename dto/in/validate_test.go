package in

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestSignUpDTOIn(t *testing.T) {
	t.Run("Valid Input", func(t *testing.T) {
		dto := SignUpDTOIn{
			PhoneNumber:  "+62-813732117646",
			Password:     "password",
			FirstName:    "John",
			MiddleName:   "Middle",
			LastName:     "Doe",
			Gender:       "Male",
			BirthDateStr: time.Now().Format("2006-01-02"),
			Preferences: SignUpPreferencesDTOIn{
				Gender: "Male",
				MinAge: 19,
				MaxAge: 25,
			},
			Passions: []string{"sports", "music", "test"},
		}

		util := validator.New()
		err := util.Struct(dto)

		assert.NoError(t, err)
	})

	t.Run("Invalid Input - Missing PhoneNumber", func(t *testing.T) {
		dto := SignUpDTOIn{
			Password:     "password",
			FirstName:    "John",
			MiddleName:   "Middle",
			LastName:     "Doe",
			Gender:       "Male",
			BirthDateStr: time.Now().Format("2006-01-02"),
			Preferences: SignUpPreferencesDTOIn{
				Gender: "Male",
				MinAge: 19,
				MaxAge: 25,
			},
			Passions: []string{"sports", "music", "test"},
		}

		util := validator.New()
		err := util.Struct(dto)

		assert.Error(t, err)
	})

	t.Run("Invalid Input - Invalid BirthDate", func(t *testing.T) {
		dto := SignUpDTOIn{
			PhoneNumber:  "+62-813732117646",
			Password:     "password",
			FirstName:    "John",
			MiddleName:   "Middle",
			LastName:     "Doe",
			Gender:       "Male",
			BirthDateStr: "invalid-date",
			Preferences: SignUpPreferencesDTOIn{
				Gender: "Male",
				MinAge: 19,
				MaxAge: 25,
			},
			Passions: []string{"sports", "music", "test"},
		}

		util := validator.New()
		err := util.Struct(dto)

		assert.Error(t, err)
	})
}

func TestLoginDTOIn(t *testing.T) {
	t.Run("Valid Input", func(t *testing.T) {
		dto := LoginDTOIn{
			PhoneNumber: "+62-813732117646",
			Password:    "password",
		}

		util := validator.New()
		err := util.Struct(dto)

		assert.NoError(t, err)
	})

	t.Run("Invalid Input - Missing PhoneNumber", func(t *testing.T) {
		dto := LoginDTOIn{
			Password: "password",
		}

		util := validator.New()
		err := util.Struct(dto)

		assert.Error(t, err)
	})

	t.Run("Invalid Input - Missing Password", func(t *testing.T) {
		dto := LoginDTOIn{
			PhoneNumber: "+62-813732117646",
		}

		util := validator.New()
		err := util.Struct(dto)

		assert.Error(t, err)
	})

	t.Run("Invalid Input - Missing Both Fields", func(t *testing.T) {
		dto := LoginDTOIn{}

		util := validator.New()
		err := util.Struct(dto)

		assert.Error(t, err)
	})
}
