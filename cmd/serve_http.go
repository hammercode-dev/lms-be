package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/hammer-code/lms-be/app"
	"github.com/hammer-code/lms-be/config"
	_ "github.com/hammer-code/lms-be/docs"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"

	// _ "swagger-mux/docs"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/swaggo/swag"
)

var serveHttpCmd = &cobra.Command{
	Use:   "http",
	Short: "launches an HTTP server",
	Long:  "the serveHttp command initiates an HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		// load add package serve http here
		ctx := context.Background()

		cfg := config.GetConfig()

		app := app.InitApp(cfg)

		// route
		router := registerHandler(app)
		router.Use(app.Middleware.LogMiddleware)

		// build cors
		muxCorsWithRouter := cors.AllowAll().Handler(router)

		srv := &http.Server{
			Addr:    cfg.APP_PORT,
			Handler: muxCorsWithRouter,
		}

		go func() {
			done := make(chan os.Signal, 1)

			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			<-done
			logrus.Info("svr.Shutdown: start shutdown")
			if err := srv.Shutdown(ctx); err == context.DeadlineExceeded {
				logrus.Error("svr.Shutdown: context deadline exceeded", err)
			}
			logrus.Info("svr.Shutdown: shutdown success")
		}()

		logrus.Info(fmt.Sprintf("server started, running on port %s", cfg.APP_PORT))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatal("starting server failed", err)
		}
	},
}

// func LoadJSON(path string) string {
// 	jsonBytes, err := os.ReadFile(path)

// 	// jsonBytes, err := os.ReadFile("documentation/users.json")
// 	if err != nil {
// 		fmt.Println("Error reading JSON file:", err)
// 		return ""
// 	}
// 	return string(jsonBytes)
// }

func LoadJSON(path string) string {
	jsonBytes, err := os.ReadFile(path)

	// jsonBytes, err := os.ReadFile("documentation/users.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return ""
	}
	return string(jsonBytes)
}

func LoadSwagger() {
	userTemplate := LoadJSON("documentation/users.json")
	var UsersSwaggerInfo = &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate:  userTemplate,
	}
	swag.Register(UsersSwaggerInfo.InstanceName(), UsersSwaggerInfo)
}

func init() {
	// LoadSwagger()
	rootCmd.AddCommand(serveHttpCmd)

}

func health(w http.ResponseWriter, _ *http.Request) {
	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "good",
		Data:    nil,
	}, w)
}

func registerHandler(app app.App) *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/health", health)

	router.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	v1 := router.PathPrefix("/api/v1").Subrouter()
	v1.HandleFunc("/newsletters/subscribe", app.NewLetterHandler.Subscribe).Methods(http.MethodPost)

	auth := v1.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", app.UserHandler.Register).Methods(http.MethodPost)
	auth.HandleFunc("/login", app.UserHandler.Login).Methods(http.MethodPost)
	auth.HandleFunc("/forgot_password", app.UserHandler.ForgotPassword).Methods(http.MethodPost)
	auth.HandleFunc("/set_password", app.UserHandler.ResetPassword).Methods(http.MethodPut)

	public := v1.PathPrefix("/public").Subrouter()
	public.HandleFunc("/storage/{kind}/{path}", app.ImageHandler.GetStorage).Methods(http.MethodGet)
	public.HandleFunc("/events", app.EventHandler.List).Methods(http.MethodGet)
	public.HandleFunc("/events/{id}", app.EventHandler.GetEventByID).Methods(http.MethodGet)
	public.HandleFunc("/images", app.ImageHandler.UploadImage).Methods(http.MethodPost)
	public.HandleFunc("/events/registrations", app.EventHandler.RegisterEvent).Methods(http.MethodPost)
	public.HandleFunc("/events/registrations/{order_no}", app.EventHandler.RegistrationStatus).Methods(http.MethodGet)
	public.HandleFunc("/events/pay", app.EventHandler.PayEvent).Methods(http.MethodPost)

	protectedV1Route := v1.NewRoute().Subrouter()
	protectedV1Route.Use(app.Middleware.AuthMiddleware)

	protectedV1AdminRoute := v1.PathPrefix("/admin").Subrouter()
	protectedV1AdminRoute.Use(app.Middleware.AuthMiddleware)

	protectedV1Route.HandleFunc("/users", app.UserHandler.GetUsers).Methods(http.MethodGet)
	protectedV1Route.HandleFunc("/user", app.UserHandler.GetUserProfile).Methods(http.MethodGet)
	protectedV1Route.HandleFunc("/logout", app.UserHandler.Logout).Methods(http.MethodPost)

	protectedV1Route.HandleFunc("/", app.UserHandler.GetUserById).Methods(http.MethodGet)
	protectedV1Route.HandleFunc("/update", app.UserHandler.UpdateProfileUser).Methods(http.MethodPut)
	protectedV1Route.HandleFunc("/delete", app.UserHandler.DeleteUser).Methods(http.MethodDelete)

	protectedV1Route.HandleFunc("/events", app.EventHandler.CreateEvent).Methods(http.MethodPost)
	protectedV1Route.HandleFunc("/events/registrations", app.EventHandler.ListRegistration).Methods(http.MethodGet)
	protectedV1Route.HandleFunc("/events/pays", app.EventHandler.ListEventPay).Methods(http.MethodGet)
	protectedV1Route.HandleFunc("/events/pays", app.EventHandler.PayProcess).Methods(http.MethodPost)
	protectedV1Route.HandleFunc("/events/{id}", app.EventHandler.GetEventByID).Methods(http.MethodGet)
	protectedV1Route.HandleFunc("/images", app.ImageHandler.UploadImage).Methods(http.MethodPost)

	protectedV1AdminRoute.HandleFunc("/events", app.EventHandler.CreateEvent).Methods(http.MethodPost)
	protectedV1AdminRoute.HandleFunc("/events", app.EventHandler.GetEvents).Methods(http.MethodGet)
	protectedV1AdminRoute.HandleFunc("/events/{id}", app.EventHandler.DeleteEvent).Methods(http.MethodDelete)
	protectedV1AdminRoute.HandleFunc("/events/{id}", app.EventHandler.GetDetail).Methods(http.MethodGet)

	protectedV1AdminRoute.HandleFunc("/users", app.UserHandler.GetUsers).Methods(http.MethodGet)

	return router
}
