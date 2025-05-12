package domain

import (
	"time"
)

type ResetPasswordToken struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	UserID     uint64    `gorm:"not null;index"` // index untuk efisiensi query
	Token      string    `gorm:"type:varchar(100);not null"`
	ExpiryDate time.Time `gorm:"not null"`
	IsUsed     bool      `gorm:"default:false"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (ResetPasswordToken) TableName() string {
	return "resetpasswordtoken"
}
