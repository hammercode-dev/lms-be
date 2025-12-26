package repository

import (
	"github.com/hammer-code/lms-be/domain"
	"golang.org/x/net/context"
)

func (r *repository) UpdateAuthorAvatar(ctx context.Context, userID uint, avatar string) error {
	return r.db.DB(ctx).Model(&domain.Author{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"avatar": avatar,
		}).Error
}
