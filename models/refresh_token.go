package models

import "time"

type RefreshToken struct {
	ID            string    `gorm:"primaryKey"`
	AccessTokenID string    `gorm:"index:access_token_id_idx"`
	ExpiredAt     time.Time `json:"expired_at"`

	AccessToken AccessToken `gorm:"foreignKey:AccessTokenID;constraint:OnDelete:CASCADE;"`
}
