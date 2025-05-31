package models

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	Email     string `gorm:"unique;index:user_email_idx"`
	Password  string `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
