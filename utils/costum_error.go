package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/ngelog"
)

const (
	ErrDuplicateEmail 	 = "\"uni_users_email\" (SQLSTATE 23505)"
	ErrWrongPassword  	 = "password"
	ErrNotFoundSQL       = "sql: no rows in result set"
)

type CustomHttpError struct {
	Code int
	Message string
	OriginError error
}

func (e *CustomHttpError) Error() string {
	if e.OriginError != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.OriginError.Error())
	}
	return e.Message
}

func NewBadRequestError(ctx context.Context, msg string, err error) *CustomHttpError {
	ngelog.Error(ctx, msg, err)
	return &CustomHttpError{
		Code:    http.StatusBadRequest,
		Message: msg,
		OriginError: err,
	}
}

func NewUnauthorizedError(ctx context.Context, msg string, err error) *CustomHttpError {
	ngelog.Error(ctx, msg, err)
	return &CustomHttpError{
		Code:    http.StatusUnauthorized,
		Message: msg,
		OriginError: err,
	}
}

func NewInternalServerError(ctx context.Context, err error) *CustomHttpError {
	ngelog.Error(ctx, "Internal server error", err)
	return &CustomHttpError{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
		OriginError: err,
	}
}

func CheckError(errStr, containsStr, message string, code int) (domain.HttpResponse, bool) {
	if strings.Contains(errStr, containsStr) {
		return domain.HttpResponse{
			Code:    code,
			Message: message,
		}, true
	}
	return domain.HttpResponse{}, false
}

func CustomErrorResponse(err error) domain.HttpResponse {
	if customErr, ok := err.(*CustomHttpError); ok {
		var errStr string
		if customErr.OriginError != nil {
			errStr = customErr.OriginError.Error()
		}
		resp, ok := CheckError(errStr, ErrDuplicateEmail, "User already exist", http.StatusBadRequest)
		if ok {
			return resp
		}

		resp, ok = CheckError(errStr, ErrWrongPassword, "Sorry, your password is incorrect", http.StatusBadRequest)
		if ok {
			return resp
		}

		resp, ok = CheckError(errStr, ErrNotFoundSQL, "Data not found", http.StatusNotFound)
		if ok {
			return resp
		}

		return domain.HttpResponse{
			Code:    customErr.Code,
			Message: customErr.Message,
		}
	}

	return domain.HttpResponse{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
	}

}
