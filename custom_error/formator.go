package custom_error

import (
	"strings"
	"time"

	"github.com/william22913/example-dating-apps/bundles"
	"github.com/william22913/example-dating-apps/dto/out"
)

var defaultLanguage = "en-US"
var appVersion = "1.0.0"
var defaultInternalCode = "E-5-DATEAPPS-SRV-001"
var _formator *formator

func NewErrorFormator(
	bundles bundles.Bundles,
) Formator {
	_formator = &formator{
		bundles: bundles,
	}
	return _formator
}

type formator struct {
	bundles bundles.Bundles
}

func NewErrorMessageParam(
	err error,
) *errorMessageParam {
	return &errorMessageParam{
		err: err,
	}
}

func (f *formator) DefaultLanguage(
	lang string,
) *formator {
	defaultLanguage = lang
	return f
}

func (f *formator) Version(
	version string,
) *formator {
	appVersion = version
	return f
}

func (f *formator) DefaultInternalCode(
	defaultCode string,
) *formator {
	defaultInternalCode = defaultCode
	return f
}

func (f formator) ReformatErrorMessage(
	param errorMessageParam,
) out.DefaultErrorResponse {

	result := out.DefaultErrorResponse{
		DefaultMessage: out.DefaultMessage{
			Success: false,
		},
	}

	language := defaultLanguage

	result.DefaultMessage.Header = out.Header{
		Version:   appVersion,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	switch errs := param.err.(type) {
	case *UnbundledErrorMessages:
		result.Payload = out.DefaultError{
			Status: errs.status,
			Code:   errs.code.Error(),
		}

		param := make(map[string]interface{})

		if errs.reason != "" {
			result.Payload.Message = errs.reason
			return result
		}

		if errs.function != nil {
			tempParam := errs.function(errs.param...)

			for key := range tempParam {
				param[key] = tempParam[key].Param
				val, _ := tempParam[key].Param.(string)

				if tempParam[key].IsConverted {
					param[key] = f.bundles.ReadMessageBundle("common.constanta", strings.ToUpper(val), language, nil)
				}
			}
		}

		result.Payload.Message = f.bundles.ReadMessageBundle("common.error", errs.Error(), language, param)
	default:
		result.Payload = out.DefaultError{
			Status:  500,
			Code:    defaultInternalCode,
			Message: f.bundles.ReadMessageBundle("common.error", defaultInternalCode, language, nil),
		}
	}

	return result
}
