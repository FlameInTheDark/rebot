package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/foundation/berrors"
	"github.com/FlameInTheDark/rebot/foundation/logs"
)

//RespondError renders error message into http server response
func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	bError := &berrors.BusinessError{}
	logger := r.Context().Value(logs.LoggerKey).(*zap.Logger)

	if errors.As(err, &bError) {
		if logger != nil {
			logger.Info(
				err.Error(),
				zap.String("ip", r.RemoteAddr),
				zap.String("uri", r.RequestURI),
				zap.Error(bError),
				zap.Int("error_code", bError.ErrCode),
			)
		}
		w.WriteHeader(http.StatusBadRequest)

		var msg string

		render.JSON(w, r, ClientError{
			Error: struct {
				Code int    `json:"code"`
				Msg  string `json:"message,omitempty"`
			}{
				Code: bError.Code(),
				Msg:  msg,
			},
		})
		return
	}

	if logger != nil {
		logger.Info(
			"Internal server error",
			zap.String("uri", r.RequestURI),
			zap.Error(bError),
			zap.Int("error_code", bError.ErrCode),
		)
	}

	w.WriteHeader(http.StatusInternalServerError)
	render.JSON(w, r, ClientError{
		Error: struct {
			Code int    `json:"code"`
			Msg  string `json:"message,omitempty"`
		}{
			Code: 0,
			Msg:  fmt.Sprintf("%s", err),
		},
	})
}

// ClientError is JSON response body containing business error code
type ClientError struct {
	Error struct {
		Code int    `json:"code"`
		Msg  string `json:"message,omitempty"` // if debug
	} `json:"error"`
}
