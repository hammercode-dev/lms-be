package utils

import (
	"net/http"
	"strings"

	"github.com/hammer-code/lms-be/domain"
)

type CustomHttpError struct {
	Code int
	Message string
	OriginError error
}

func (e *CustomHttpError) Error() string {
	return e.Message
}

func NewBadRequestError(msg string, err error) *CustomHttpError {
	return &CustomHttpError{
		Code:    http.StatusBadRequest,
		Message: msg,
		OriginError: err,
	}
}

func NewUnauthorizedError(msg string, err error) *CustomHttpError {
	return &CustomHttpError{
		Code:    http.StatusUnauthorized,
		Message: msg,
		OriginError: err,
	}
}

func NewInternalServerError(err error) *CustomHttpError {
	return &CustomHttpError{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
		OriginError: err,
	}
}

func CheckError(err, sub, message string, code int) (domain.HttpResponse, bool) {
	if strings.Contains(err, sub) {
		return domain.HttpResponse{
			Code:    code,
			Message: message,
		}, true
	}
	return domain.HttpResponse{}, false
}

func CustomErrorResponse(err error) domain.HttpResponse {
	if customErr, ok := err.(*CustomHttpError); ok {
		errStr := customErr.OriginError.Error()
		resp, ok := CheckError(errStr, "\"uni_users_email\" (SQLSTATE 23505)", "User already exist", 400)
		if ok {
			return resp
		}

		resp, ok = CheckError(errStr, "\"uni_logout_token\" (SQLSTATE 23505)", "You have already logged out.", 400)
		if ok {
			return resp
		}

		resp, ok = CheckError(errStr, "password", "Sorry, your password is incorrect", 400)
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
