package repository

import (
	"errors"

	"github.com/hammer-code/lms-be/domain"
	pkgDB "github.com/hammer-code/lms-be/pkg/db"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type repository struct {
	db pkgDB.DatabaseTransaction
}

// GetDetailBlogPost implements domain.BlogPostRepository.
func (r *repository) GetDetailBlogPost(ctx context.Context, slug, typeFind string, id uint) (data domain.BlogPost, err error) {
	db := r.db.DB(ctx).Preload("Author").Model(&domain.BlogPost{}).Where("is_deleted = ?", false)

	switch typeFind {
	case "slug":
		err = db.First(&data, "slug = ?", slug).Error
		if err != nil {
			logrus.Error("failed to get blog post detail: ", err)
			return data, err
		}
	case "id":
		err = db.First(&data, "id = ?", id).Error
		if err != nil {
			logrus.Error("failed to get blog post detail: ", err)
			return data, err
		}
	default:
		return domain.BlogPost{}, errors.New("invalid typeFind parameter, must be 'slug' or 'id'")
	}

	var tags []string
	if err := r.db.DB(ctx).Table("blog_post_tags").
		Select("tag").
		Where("blog_post_id = ?", data.Id).
		Pluck("tag", &tags).Error; err != nil {
		logrus.Error("failed to get tags for blog post ID ", data.Id, ": ", err)
	} else {
		data.Tags = tags
	}

	return data, nil
}

// UpdateBlogPost implements domain.BlogPostRepository.
func (r *repository) UpdateBlogPost(ctx context.Context, data domain.BlogPost, id uint) error {
	return r.db.StartTransaction(ctx, func(txCtx context.Context) error {
		if err := r.db.DB(txCtx).Model(&domain.BlogPost{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"title":        data.Title,
				"content":      data.Content,
				"excerpt":      data.Excerpt,
				"published_at": data.PublishedAt,
				"updated_at":   data.UpdatedAt,
				"category":     data.Category,
				"status":       data.Status,
			}).Error; err != nil {
			logrus.Error("failed to update blog post: ", err)
			return err
		}
		if data.Author.Avatar != "" {
			if err := r.db.DB(txCtx).Model(&domain.Author{}).
				Where("user_id = ?", data.AuthorID).
				Updates(map[string]interface{}{
					"avatar": data.Author.Avatar,
				}).Error; err != nil {
				logrus.Error("failed to update author avatar: ", err)
				return err
			}
		}

		if len(data.Tags) > 0 {
			if err := r.db.DB(txCtx).Table("blog_post_tags").
				Where("blog_post_id = ?", id).
				Delete(nil).Error; err != nil {
				logrus.Error("failed to delete old tags: ", err)
				return err
			}

			for _, tag := range data.Tags {
				blogPostTag := struct {
					BlogPostId int    `gorm:"column:blog_post_id"`
					Tag        string `gorm:"column:tag"`
				}{
					BlogPostId: int(id),
					Tag:        tag,
				}

				if err := r.db.DB(txCtx).Table("blog_post_tags").
					Create(&blogPostTag).Error; err != nil {
					logrus.Error("failed to create blog post tag: ", err)
					return err
				}
			}
		}

		return nil
	})
}

// CreateBlogPost implements domain.BlogPostRepository.
func (r *repository) CreateBlogPost(ctx context.Context, data domain.BlogPost) error {
	// Menggunakan StartTransaction yang sudah ada
	err := r.db.StartTransaction(ctx, func(txCtx context.Context) error {
		// 1. Periksa/Buat Author jika belum ada
		var authorExists int64
		if err := r.db.DB(txCtx).Model(&domain.Author{}).
			Preload("Author").
			Where("user_id = ?", data.Author.UserId).
			Count(&authorExists).Error; err != nil {
			logrus.Error("failed to check author existence: ", err)
			return err
		}

		if authorExists == 0 {
			// Insert author
			if err := r.db.DB(txCtx).Create(&data.Author).Error; err != nil {
				logrus.Error("failed to create author: ", err)
				return err
			}
		}

		data.UpdatedAt = nil
		// Set AuthorID untuk relasi
		data.AuthorID = data.Author.UserId

		// 2. Insert Blog Post
		if err := r.db.DB(txCtx).Omit("updated_at").Create(&data).Error; err != nil {
			logrus.Error("failed to create blog post: ", err)
			return err
		}

		// 3. Insert Tags jika ada
		if len(data.Tags) > 0 {
			for _, tag := range data.Tags {
				blogPostTag := struct {
					BlogPostId int    `gorm:"column:blog_post_id"`
					Tag        string `gorm:"column:tag"`
				}{
					BlogPostId: data.Id,
					Tag:        tag,
				}

				if err := r.db.DB(txCtx).Table("blog_post_tags").Create(&blogPostTag).Error; err != nil {
					logrus.Error("failed to create blog post tag: ", err)
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// DeleteBlogPost implements domain.BlogPostRepository.
func (r *repository) DeleteBlogPost(ctx context.Context, id uint) error {
	db := r.db.DB(ctx).Model(&domain.BlogPost{})

	// Perform soft delete by updating is_deleted field
	result := db.Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
	})

	if result.Error != nil {
		logrus.Error("failed to soft delete blog post: ", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		logrus.Warn("no blog post found to delete with id: ", id)
		return errors.New("blog post not found")
	}
	return nil
}

// GetAllBlogPosts implements domain.BlogPostRepository.
func (r *repository) GetAllBlogPosts(ctx context.Context, pagination domain.FilterPagination) ([]domain.BlogPost, int, error) {
	var data []domain.BlogPost
	var totalCount int64

	if err := r.db.DB(ctx).Model(&domain.BlogPost{}).Where("is_deleted = ?", false).Count(&totalCount).Error; err != nil {
		logrus.Error("failed to count blog posts: ", err)
		return nil, 0, err
	}

	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	orderBy := pagination.GetOrderBy()

	query := r.db.DB(ctx).Preload("Author").Where("is_deleted = ?", false)

	if orderBy != "" {
		query = query.Order(orderBy)
	} else {
		query = query.Order("id DESC")
	}

	err := query.Limit(limit).Offset(offset).Find(&data).Error
	if err != nil {
		logrus.Error("failed to get all blog posts: ", err)
		return nil, 0, err
	}

	for i := range data {
		var tags []string
		if err := r.db.DB(ctx).Table("blog_post_tags").
			Select("tag").
			Where("blog_post_id = ?", data[i].Id).
			Pluck("tag", &tags).Error; err != nil {
			logrus.Error("failed to get tags for blog post ID ", data[i].Id, ": ", err)
		} else {
			data[i].Tags = tags
		}
	}

	return data, int(totalCount), nil
}

func NewRepository(db pkgDB.DatabaseTransaction) domain.BlogPostRepository {
	return &repository{
		db: db,
	}
}
