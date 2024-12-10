package custom_endpoint

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/william22913/example-dating-apps/constanta"
	"github.com/william22913/example-dating-apps/custom_context"
	"github.com/william22913/example-dating-apps/custom_error"
	"github.com/william22913/example-dating-apps/dto/out"
)

func (g HTTPController) WrapService(
	param *wrapServiceParam,
) func(
	http.ResponseWriter,
	*http.Request,
) {
	return func(
		rw http.ResponseWriter,
		r *http.Request,
	) {
		var ctx *custom_context.ContextModel
		var err error
		var payload interface{}
		var header map[string]string

		timeNow := time.Now()

		defer func() {
			g.response(
				ctx,
				r,
				rw,
				timeNow,
				header,
				payload,
				err,
			)
		}()

		ctx = g.getContext(r)

		if err != nil {
			err = custom_error.ErrReadBody
			return
		}

		ctx.ClientAccess.Path = r.URL.Path

		// Validate Token
		_, headers := g.Converter(r)

		ctx.ClientAccess.Timestamp = time.Now()

		if param.cv != nil {
			err = param.cv(ctx.ToContext(), headers)
			if err != nil {
				return
			}

			c, valid := r.Context().Value(
				constanta.ApplicationContextConstanta,
			).(*custom_context.ContextModel)

			if valid {
				ctx = c
			}
		} else {
			ctx = custom_context.NewContextModel()
		}

		request_id := ReadHeader(r, "X-Request-Id")
		if request_id == "" {
			ctx.ClientAccess.RequestID = GetUUID()
		} else {
			ctx.ClientAccess.RequestID = request_id
		}

		var dto interface{}

		dto, err = g.readBody(
			ctx,
			param,
			r,
		)

		if err != nil {
			return
		}

		header, payload, err = param.serve(
			ctx,
			dto,
		)
	}
}

func (g HTTPController) response(
	ctx *custom_context.ContextModel,
	r *http.Request,
	rw http.ResponseWriter,
	timeNow time.Time,
	header map[string]string,
	payload interface{},
	err error,
) {
	var errs out.DefaultErrorResponse
	var usedErr = err
	statusCode := 200
	var length int

	defer func() {

		usedPath, _ := mux.CurrentRoute(r).GetPathTemplate()
		if usedPath == "" {
			usedPath = r.URL.Path
		}

		msg := "Api Called"
		if errs.Payload.Message != "" {
			msg = errs.Payload.Message
		}

		logEvent := log.Info().
			Str("action", "middleware.api.call").
			Str("url", usedPath).
			Str("method", r.Method).
			Int("status", statusCode).
			Int("out_size", length).
			Int("process_time", int(time.Since(timeNow).Seconds())).
			Str("request_id", ctx.ClientAccess.RequestID)

		if usedErr != nil {
			logEvent = logEvent.Err(usedErr)
		}

		logEvent.Msg(msg)

	}()

	rw.Header().Set(constanta.RequestIDConstanta, ctx.ClientAccess.RequestID)

	for key := range header {
		headersList := rw.Header().Get("Access-Control-Allow-Headers") + ", " + strings.ToLower(key)
		rw.Header().Set("Access-Control-Allow-Headers", headersList)
		rw.Header().Set("Access-Control-Expose-Headers", headersList)
		rw.Header().Add(key, header[key])
	}

	var data []byte
	if err != nil {
		errs = g.formator.ReformatErrorMessage(
			*custom_error.NewErrorMessageParam(
				err,
			),
		)

		errs.Header.RequestID = ctx.ClientAccess.RequestID
		statusCode = errs.Payload.Status
		payload = errs

	} else {
		payload = dtoConverter.GetFunction()(ctx, payload)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)

	if data == nil {
		data, err = json.Marshal(payload)
		if err != nil {
			log.Error().
				Err(err).
				Caller().
				Msg("Error on Marshaling Message")
			return
		}

		length, err = rw.Write(data)
		if err != nil {
			log.Error().
				Err(err).
				Caller().
				Msg("Error on Writing Message")
			return
		}
	}

}
