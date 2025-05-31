package entities

import "time"

type AccessToken struct {
	ID        string
	UserID    uint
	ExpiredAt time.Time
}

type RefreshToken struct {
	ID            string
	AccessTokenID string
	ExpiredAt     time.Time
}
