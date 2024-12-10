package sign_up

import (
	"github.com/william22913/example-dating-apps/custom_context"
)

type SignUpService interface {
	SignUpUser(
		*custom_context.ContextModel,
		interface{},
	) (
		map[string]string,
		interface{},
		error,
	)

	GetDTO() interface{}
}
