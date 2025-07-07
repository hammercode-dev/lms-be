package blog_post

import (
	blog_post_handler "github.com/hammer-code/lms-be/app/blog_post/delivery/http"
	blog_post_repo "github.com/hammer-code/lms-be/app/blog_post/repository"
	blog_post_usecase "github.com/hammer-code/lms-be/app/blog_post/usecase"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/db"
	"github.com/hammer-code/lms-be/pkg/jwt"
)

func InitRepository(db db.DatabaseTransaction) domain.BlogPostRepository {
	return blog_post_repo.NewRepository(db)
}

func InitUseCase(repository domain.BlogPostRepository, jwt jwt.JWT) domain.BlogPostUsecase {
	return blog_post_usecase.NewUsecase(repository, jwt)
}

func InitHandler(usecase domain.BlogPostUsecase) domain.BlogPostHandler {
	return blog_post_handler.NewHandler(usecase)
}
