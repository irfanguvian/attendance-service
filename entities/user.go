package entities

import "time"

type User struct {
	ID        uint
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
