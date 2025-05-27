package custresp

import (
	"database/sql"
	"errors"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response/rest"
	"github.com/nocturna-ta/result/pkg/constants"
	"github.com/nocturna-ta/result/pkg/constants/errorcode"
	"net/http"
)

func CustomErrorResponse(err error) (*rest.JSONResponse, error) {
	resp := rest.NewJSONResponse()
	if err == nil {
		return resp, nil
	}

	var e *custerr.ErrChain
	switch {
	case errors.As(err, &e):
		errCause := e.Cause
		if errCause != nil {
			err = errCause
		}

		message := e.Message
		if message == constants.EmptyString {
			message = err.Error()
		}

		resp.SetCode(getErrorCode(e))
		resp.Error = &rest.ErrorResponse{
			ErrorCode:    e.Code,
			ErrorMessage: message,
		}
		return resp, nil
	case errors.Is(err, dao.ErrNoResult), errors.Is(err, sql.ErrNoRows):
		resp.SetCode(http.StatusBadRequest)
		resp.Error = &rest.ErrorResponse{
			ErrorCode:    errorcode.NotFound.Code,
			ErrorMessage: errorcode.NotFound.Message,
		}
		return resp, nil
	default:
		return resp.SetError(err).SetMessage(err.Error()), nil
	}

}

func getErrorCode(err *custerr.ErrChain) int {
	switch {
	case errors.Is(err.Type, ErrTooManyRequest):
		return http.StatusTooManyRequests
	case errors.Is(err.Type, ErrRequestTooEarly):
		return http.StatusTooEarly
	case errors.Is(err.Type, ErrInvalidRequest):
		return http.StatusNotAcceptable
	default:
		return rest.GetErrorCode(err)
	}
}
