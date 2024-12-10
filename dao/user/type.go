package user

import (
	"github.com/william22913/example-dating-apps/repository"
)

type UserDAO interface {
	CheckIsPhoneExist(
		repository.User,
	) (
		repository.User,
		error,
	)

	InsertUserData(
		repository.CompletedUserData,
	) error

	GetUserDataForLogin(
		data repository.User,
	) (
		repository.CompletedUserData,
		error,
	)
}
