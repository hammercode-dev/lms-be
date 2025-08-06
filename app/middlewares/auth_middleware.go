package middlewares

import (
	"context"
	"net/http"
	"strconv"

	"github.com/hammer-code/lms-be/domain"
	contextkey "github.com/hammer-code/lms-be/pkg/context_key"
	"github.com/hammer-code/lms-be/pkg/ngelog"
	"github.com/hammer-code/lms-be/utils"
)

func (m *Middleware) AuthMiddleware(allowedRole string) domain.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx, span := tracer.Start(request.Context(), "auth middleware")
			defer span.End()

			token := utils.ExtractBearerToken(request)
			if len(*token) < 5 {
				ngelog.Error(ctx, "failed to extract bearer token", nil)
				utils.Response(domain.HttpResponse{
					Code:    401,
					Message: "Unauthorized",
					Data:    nil,
				}, writer)
				return
			}

			verifyToken, err := m.Jwt.VerifyToken(*token)
			if err != nil {
				ngelog.Error(ctx, "failed to verify token", err)
				utils.Response(domain.HttpResponse{
					Code:    500,
					Message: "failed to verify token",
					Data:    nil,
				}, writer)
				return
			}

			logoutToken, err := m.UserRepo.GetToken(request.Context(), *token)
			if err != nil {
				ngelog.Error(ctx, "failed to get token", err)
				utils.Response(domain.HttpResponse{
					Code:    401,
					Message: "Unauthorized",
					Data:    nil,
				}, writer)
				return
			}
			if logoutToken.Status == 0 {
				ngelog.Error(ctx, "unauthorized", nil)
				utils.Response(domain.HttpResponse{
					Code:    401,
					Message: "Unauthorized",
					Data:    nil,
				}, writer)
				return
			}

			user, err := m.UserRepo.FindByEmail(request.Context(), verifyToken.Email)
			if err != nil {
				ngelog.Error(ctx, "failed to find by email", err)
				utils.Response(domain.HttpResponse{
					Code:    401,
					Message: "Unauthorized",
					Data:    nil,
				}, writer)
				return
			}

			if user.Role != allowedRole {
				ngelog.Error(ctx, "role is not the role", nil)
				utils.Response(domain.HttpResponse{
					Code:    401,
					Message: "Unauthorized",
					Data:    nil,
				}, writer)
				return
			}

			writer.Header().Set("x-user-id", strconv.Itoa(user.ID))
			writer.Header().Set("x-username", user.Username)
			
			ctxUser := context.WithValue(request.Context(), contextkey.UserKey, user.ID)
			request = request.WithContext(ctxUser)

			next.ServeHTTP(writer, request)
		})
	}
}
