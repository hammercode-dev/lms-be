package usecase

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type usecase struct {
	repo domain.BlogPostRepository
	jwt  jwt.JWT
}

// CreateBlogPost implements domain.BlogPostUsecase.
func (uc *usecase) CreateBlogPost(ctx context.Context, data domain.BlogPost, token string) (domain.BlogPost, error) {

	jwtData, err := uc.jwt.VerifyToken(token)
	if err != nil {
		logrus.Error("failed to verify token: ", err)
		return domain.BlogPost{}, err
	}

	slugBytes := make([]byte, 32)
	if _, err := rand.Read(slugBytes); err != nil {
		return domain.BlogPost{}, err
	}

	data.Author.UserId = jwtData.ID
	data.Author.Name = jwtData.UserName
	data.Slug = hex.EncodeToString(slugBytes)

	blogPost, err := uc.repo.CreateBlogPost(ctx, data)
	if err != nil {
		logrus.Error("failed to create blog post: ", err)
		return domain.BlogPost{}, err

	}

	return blogPost, nil

}

// DeleteBlogPost implements domain.BlogPostUsecase.
func (uc *usecase) DeleteBlogPost(ctx context.Context, id uint) error {
	if err := uc.repo.DeleteBlogPost(ctx, id); err != nil {
		logrus.Error("failed to delete blog post detail: ", err)
		return err
	}

	return nil
}

// GetAllBlogPosts implements domain.BlogPostUsecase.
func (uc *usecase) GetAllBlogPosts(ctx context.Context) ([]domain.BlogPost, error) {
	blogPosts, err := uc.repo.GetAllBlogPosts(ctx)
	if err != nil {
		logrus.Error("failed to get all blog posts: ", err)
		return nil, err
	}
	return blogPosts, nil
}

// GetDetailBlogPost implements domain.BlogPostUsecase.
func (uc *usecase) GetDetailBlogPost(ctx context.Context, slug string, id uint) (domain.BlogPost, error) {
	typeFind := "slug"
	if slug == "" {
		typeFind = "id"
	}
	blogPost, err := uc.repo.GetDetailBlogPost(ctx, slug, typeFind, id)
	if err != nil {
		logrus.Error("failed to get blog post detail: ", err)
		return domain.BlogPost{}, err
	}
	return blogPost, nil
}

// UpdateBlogPost implements domain.BlogPostUsecase.
func (uc *usecase) UpdateBlogPost(ctx context.Context, data domain.BlogPost, id uint) error {
	err := uc.repo.UpdateBlogPost(ctx, data, id)
	if err != nil {
		logrus.Error("failed to update blog post: ", err)
		return err
	}
	return nil
}

var (
	usec *usecase
)

func NewUsecase(repo domain.BlogPostRepository, jwt jwt.JWT) domain.BlogPostUsecase {
	if usec == nil {
		usec = &usecase{
			repo: repo,
			jwt:  jwt,
		}
	}
	return usec
}
