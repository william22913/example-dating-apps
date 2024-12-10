package custom_error

import "errors"

var ErrUnauthorized = NewUnBundledErrorMessages(401, errors.New("E-1-CMD-AUT-001"), nil)
var ErrExpiredToken = NewUnBundledErrorMessages(401, errors.New("E-1-CMD-AUT-002"), nil)
var ErrReadBody = NewUnBundledErrorMessages(400, errors.New("E-1-CMD-CTR-001"), nil)
var ErrMarshalingBody = NewUnBundledErrorMessages(400, errors.New("E-1-CMD-CTR-002"), nil)
var ErrValidationBody = NewUnBundledErrorMessages(400, errors.New("E-1-CMD-CTR-003"), errValidationBodyConverter)
var ErrDataUsed = NewUnBundledErrorMessages(400, errors.New("E-1-CMD-SRV-001"), errValidationDataUsed)

var errValidationBodyConverter = func(value ...interface{}) map[string]ErrorParam {
	result := make(map[string]ErrorParam)
	result["FieldName"] = ErrorParam{value[0], true}
	result["FieldTag"] = ErrorParam{value[1], true}
	return result
}

var errValidationDataUsed = func(value ...interface{}) map[string]ErrorParam {
	result := make(map[string]ErrorParam)
	result["FieldName"] = ErrorParam{value[0], true}
	return result
}
