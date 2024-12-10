package custom_endpoint

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/william22913/example-dating-apps/custom_context"
	"github.com/william22913/example-dating-apps/custom_error"
)

func NewHTTPController(
	dtoValidator *validator.Validate,
) *HTTPController {
	validator := &HTTPController{
		version:      "1.0.0",
		dtoValidator: dtoValidator,
	}

	dtoConverter = &defaultResponseDTO{}

	knownHeader = map[string]string{
		"Authorization": "",
		"X-Request-Id":  "",
	}

	return validator
}

func (g *HTTPController) AddKnownHeader(param ...string) {
	for i := 0; i < len(param); i++ {
		knownHeader[param[i]] = ""
	}
}

func (g *HTTPController) Router(router *mux.Router) *HTTPController {
	g.router = router
	return g
}

func (g *HTTPController) Formator(formator custom_error.Formator) *HTTPController {
	g.formator = formator
	return g
}

func (g *HTTPController) Version(version string) *HTTPController {
	g.version = version
	dtoConverter.Version(version)
	return g
}

type HTTPController struct {
	router       *mux.Router
	formator     custom_error.Formator
	version      string
	dtoValidator *validator.Validate
}

var dtoConverter DTOResponseConverter

type DTOResponseConverter interface {
	GetFunction() func(ctx *custom_context.ContextModel, payload interface{}) interface{}
	Version(version string)
}

func (g *HTTPController) ModifyDTOResponse(
	f DTOResponseConverter,
) {
	dtoConverter = f
}

func (g HTTPController) HandleFunc(
	param *handleFuncParam,
) {
	g.router.HandleFunc(param.path, param.f).Methods(param.method...)
}

var knownHeader map[string]string

func (g HTTPController) Converter(
	r *http.Request,
) (
	context.Context,
	map[string]string,
) {
	result := make(map[string]string)

	for keys := range knownHeader {
		value := ReadHeader(r, keys)
		if value != "" {
			result[keys] = value
		}
	}

	if r.Context() != nil {
		r = r.WithContext(context.Background())
	}

	return r.Context(), result
}
