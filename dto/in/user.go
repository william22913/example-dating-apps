package in

import "time"

type SignUpDTOIn struct {
	PhoneNumber  string                 `json:"phone_number" validate:"required,min=10,max=20"`
	Password     string                 `json:"password" validate:"required,min=8,max=20"`
	FirstName    string                 `json:"first_name" validate:"required,min=1,max=50"`
	MiddleName   string                 `json:"middle_name" validate:"max=50"`
	LastName     string                 `json:"last_name" validate:"required,min=1,max=50"`
	BirthDateStr string                 `json:"birth_date" validate:"required,datetime=2006-01-02"`
	Gender       string                 `json:"gender" validate:"required,oneof=Male Female Any"`
	Preferences  SignUpPreferencesDTOIn `json:"preferences" validate:"required"`
	Passions     []string               `json:"passions" validate:"required,min=3,max=50"`
	BirthDate    time.Time
}

type SignUpPreferencesDTOIn struct {
	Gender string `json:"gender" validate:"required,oneof=Male Female Any"`
	MinAge int    `json:"min_age" validate:"min=18,max=99"`
	MaxAge int    `json:"max_age" validate:"min=18,max=99"`
}

type LoginDTOIn struct {
	PhoneNumber string `json:"phone_number" validate:"required,min=10,max=20"`
	Password    string `json:"password" validate:"required,min=8,max=20"`
}
