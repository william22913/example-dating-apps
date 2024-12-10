package custom_error

import (
	"strings"

	"github.com/william22913/example-dating-apps/dto/out"
)

type ErrorParam struct {
	Param       interface{}
	IsConverted bool
}

type Converter func(...interface{}) map[string]ErrorParam

type UnbundledErrorMessages struct {
	status   int
	code     error
	reason   string
	function Converter
	param    []interface{}
}

type errorMessageParam struct {
	err error
}

func NewUnBundledErrorMessages(
	status int,
	code error,
	f Converter,
) *UnbundledErrorMessages {
	return &UnbundledErrorMessages{
		status:   status,
		code:     code,
		function: f,
		param: []interface{}{
			"UNDEFINED",
			"UNDEFINED",
			"UNDEFINED",
			"UNDEFINED",
			"UNDEFINED",
		},
	}
}

func (e UnbundledErrorMessages) Error() string {
	if _formator == nil {
		return e.code.Error()
	} else {
		var result string
		language := defaultLanguage

		param := make(map[string]interface{})
		if e.function != nil {
			tempParam := e.function(e.param...)

			for key := range tempParam {
				param[key] = tempParam[key].Param
				val, _ := tempParam[key].Param.(string)

				if tempParam[key].IsConverted {
					param[key] = _formator.bundles.ReadMessageBundle("common.constanta", strings.ToUpper(val), language, nil)
				}
			}
		}

		result = _formator.bundles.ReadMessageBundle("common.error", e.code.Error(), language, param)

		if e.reason != "" {
			result += " ,Reason : " + e.reason
		}

		return result
	}
}

func (e *UnbundledErrorMessages) Reason(reason string) *UnbundledErrorMessages {
	e.reason = reason
	return e
}

func (e *UnbundledErrorMessages) Param(param ...interface{}) *UnbundledErrorMessages {
	e.param = param
	return e
}

type Formator interface {
	ReformatErrorMessage(
		param errorMessageParam,
	) out.DefaultErrorResponse

	DefaultLanguage(
		lang string,
	) *formator

	Version(
		version string,
	) *formator

	DefaultInternalCode(
		defaultCode string,
	) *formator
}
