package middlewares

import (
	"net/http"
	"strconv"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		token := utils.ExtractBearerToken(request)
		if len(*token) < 5 {
			utils.Response(domain.HttpResponse{
				Code:    401,
				Message: "Forbidden",
				Data:    nil,
			}, writer)
			return
		}

		verifyToken, err := m.Jwt.VerifyToken(*token)
		if err != nil {
			utils.Response(domain.HttpResponse{
				Code:    500,
				Message: err.Error(),
				Data:    nil,
			}, writer)
			return
		}

		user, err := m.UserRepo.FindByEmail(request.Context(), verifyToken.Email)
		if err != nil {
			utils.Response(domain.HttpResponse{
				Code:    401,
				Message: "Forbidden",
				Data:    nil,
			}, writer)
			return
		}

		writer.Header().Set("x-user-id", strconv.Itoa(user.ID))
		writer.Header().Set("x-username", user.Username)

		next.ServeHTTP(writer, request)
	})
}
