package repository

import (
	"golang.org/x/net/context"
)

func (r *repository) DeleteTagsByBlogPostID(ctx context.Context, blogPostID uint) error {
	return r.db.DB(ctx).Table("blog_post_tags").
		Where("blog_post_id = ?", blogPostID).
		Delete(nil).Error
}
