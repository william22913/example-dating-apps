package custom_endpoint

import (
	"context"
	"net/http"
	"time"

	"github.com/william22913/example-dating-apps/custom_context"
	"github.com/william22913/example-dating-apps/dto/out"
)

type Services interface {
	GetDTO() interface{}
}

type HTTPServerAccess func(
	ctx context.Context,
	r *http.Request,
) error

type FunctionServe func(
	*custom_context.ContextModel,
	interface{},
) (
	map[string]string,
	interface{},
	error,
)

type HTTPControllerParam struct {
	Ctx context.Context
	R   *http.Request
}

func NewWarpServiceParam(
	service Services,
	serve FunctionServe,
	cv ServerAccessValidator,
) *wrapServiceParam {
	return &wrapServiceParam{
		cv:      cv,
		service: service,
		serve:   serve,
	}
}

type wrapServiceParam struct {
	serve   FunctionServe
	service Services
	cv      ServerAccessValidator
}

type defaultResponseDTO struct {
	version string
}

func (d *defaultResponseDTO) Version(param string) {
	d.version = param
}

func (d defaultResponseDTO) GetFunction() func(ctx *custom_context.ContextModel, payload interface{}) interface{} {
	return func(ctx *custom_context.ContextModel, payload interface{}) interface{} {
		result := out.DefaultResponse{
			DefaultMessage: out.DefaultMessage{
				Success: true,
			},
		}

		result.DefaultMessage.Header = out.Header{
			Timestamp: time.Now().Format(time.RFC3339),
			Version:   d.version,
			RequestID: ctx.ClientAccess.RequestID,
		}

		result.Payload = payload
		payload = result

		return result
	}
}
