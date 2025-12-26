package repository

import "golang.org/x/net/context"

func (r *repository) GetTagsByBlogPostID(ctx context.Context, blogPostID uint) (tags []string, err error) {
	if err := r.db.DB(ctx).Table("blog_post_tags").
		Select("tag").
		Where("blog_post_id = ?", blogPostID).
		Pluck("tag", &tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
