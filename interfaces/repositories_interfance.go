package interfaces

import (
	"github.com/irfanguvian/attendance-service/models"
)

type Repositories struct {
	UserRepository UserRepository
}

type UserRepository interface {
	CreateUser(email, password string) (uint, error)
	CreateAccessToken(userID uint, accessID string) error
	CreateRefreshToken(accessID string, refreshID string) error

	DeleteAccessTokenByUserID(userID uint) error

	GetUserByEmail(email string) (*models.User, error)
	GetRefreshTokenByAccessID(accessID string) (string, error)

	GetUserByAccessID(accessID string) (*models.User, error)
}
