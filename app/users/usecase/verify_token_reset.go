package usecase

import "context"

func (u *usecase) VerifyPasswordResetToken(ctx context.Context, token string) error {
	_, err := u.jwt.VerifyToken(token)
	if err != nil {
		return err
	}
	return nil
}
