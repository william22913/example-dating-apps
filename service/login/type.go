package login

import (
	"github.com/william22913/example-dating-apps/custom_context"
)

type LoginService interface {
	Login(
		*custom_context.ContextModel,
		interface{},
	) (
		map[string]string,
		interface{},
		error,
	)

	GetDTO() interface{}
}
