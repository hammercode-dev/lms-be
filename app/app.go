package app

import (
	blogPost "github.com/hammer-code/lms-be/app/blog_post"
	"github.com/hammer-code/lms-be/app/middlewares"
	newsletters "github.com/hammer-code/lms-be/app/newsletters"
	users "github.com/hammer-code/lms-be/app/users"
	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/domain"
	pkgDB "github.com/hammer-code/lms-be/pkg/db"
	"github.com/hammer-code/lms-be/pkg/jwt"
	"gorm.io/driver/postgres"

	events "github.com/hammer-code/lms-be/app/events"
	images "github.com/hammer-code/lms-be/app/images"
)

type App struct {
	Middleware       domain.Middleware
	UserHandler      domain.UserHandler
	NewLetterHandler domain.NewslettterHandler
	EventHandler     domain.EventHandler
	ImageHandler     domain.ImageHandler
	BlogPostHandler  domain.BlogPostHandler
}

func InitApp(
	cfg config.Config,
) App {

	db := config.GetDatabase(postgres.Dialector{
		Config: &postgres.Config{
			DSN: cfg.DB_POSTGRES_DSN,
		}})

	dbTx := pkgDB.NewDBTransaction(db)
	jwtInstance := jwt.NewJwt(cfg.JWT_SECRET_KEY)

	// repository
	userRepo := users.InitRepository(dbTx)
	newsletterRepo := newsletters.InitRepository(dbTx)
	eventRepo := events.InitRepository(dbTx)
	imgRepo := images.InitRepository(dbTx)
	blogPostRepo := blogPost.InitRepository(dbTx)

	// Middlewares
	middleware := middlewares.InitMiddleware(jwtInstance, userRepo)

	// usecase
	userUsecase := users.InitUsecase(cfg, userRepo, dbTx, jwtInstance)
	newsletterUC := newsletters.InitUsecase(cfg, newsletterRepo, dbTx, jwt.NewJwt(cfg.JWT_SECRET_KEY))
	eventUC := events.InitUsecase(cfg, eventRepo, imgRepo, dbTx, jwtInstance)
	imgUc := images.InitUsecase(imgRepo, dbTx)
	blogPostUc := blogPost.InitUseCase(blogPostRepo, jwtInstance)

	// handler
	userHandler := users.InitHandler(userUsecase)
	newsletterHandler := newsletters.InitHandler(newsletterUC, middleware)
	eventHandler := events.InitHandler(eventUC)
	ImageHandler := images.InitHandler(imgUc)
	blogPostHandler := blogPost.InitHandler(blogPostUc)

	return App{
		UserHandler:      userHandler,
		NewLetterHandler: newsletterHandler,
		Middleware:       middleware,
		EventHandler:     eventHandler,
		ImageHandler:     ImageHandler,
		BlogPostHandler:  blogPostHandler,
	}
}
