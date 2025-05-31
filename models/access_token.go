package models

import "time"

type AccessToken struct {
	ID        string    `gorm:"primaryKey"`
	UserID    uint      `gorm:"index:user_id_idx"`
	ExpiredAt time.Time `json:"expired_at"`

	User User `gorm:"foreignKey:UserID"`
}
