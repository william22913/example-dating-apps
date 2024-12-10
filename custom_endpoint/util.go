package custom_endpoint

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/william22913/example-dating-apps/constanta"
	"github.com/william22913/example-dating-apps/custom_context"
	"github.com/william22913/example-dating-apps/custom_error"
)

func (g HTTPController) getContext(
	r *http.Request,
) *custom_context.ContextModel {

	var ctx *custom_context.ContextModel

	rCtx := r.Context().Value(constanta.ApplicationContextConstanta)
	if rCtx == nil {

		ctx = custom_context.NewContextModel()
		context := context.WithValue(
			r.Context(),
			constanta.ApplicationContextConstanta,
			ctx,
		)

		r = r.WithContext(context)
	} else {
		ctx = rCtx.(*custom_context.ContextModel)
	}

	requestID := ReadHeader(r, constanta.RequestIDConstanta)
	if requestID == "" {
		requestID = GetUUID()
		r.Header.Set(constanta.RequestIDConstanta, requestID)
	}

	return ctx

}

func GetUUID() (output string) {
	UUID, _ := uuid.NewRandom()
	output = UUID.String()
	output = strings.Replace(output, "-", "", -1)
	return
}

func (g HTTPController) readBody(
	ctx *custom_context.ContextModel,
	param *wrapServiceParam,
	request *http.Request,
) (
	dto interface{},
	err error,
) {
	var stringBody string

	if request.Method != "GET" {

		dto = param.service.GetDTO()

		stringBody, _, err = ReadBody(request)
		if err != nil {
			return nil, custom_error.ErrReadBody
		}

		err = json.Unmarshal([]byte(stringBody), dto)
		if err != nil {
			return nil, custom_error.ErrMarshalingBody
		}

		err = g.dtoValidator.Struct(dto)
		if err != nil {
			for _, errs := range err.(validator.ValidationErrors) {
				return nil, custom_error.ErrValidationBody.Param(errs.Field(), errs.Tag())
			}
		}

	}

	return
}

func ReadBody(request *http.Request) (output string, bodySize int, err error) {
	byteBody, err := io.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		return "", 0, errors.New("BODY_INVALID")
	}
	return string(byteBody), len(byteBody), nil
}

func ReadHeader(request *http.Request, headerName string) (result string) {
	result = request.Header.Get(headerName)
	if result == "" {
		result = request.Header.Get(strings.ToLower(headerName))
	}
	return result
}
