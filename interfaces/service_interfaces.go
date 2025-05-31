package interfaces

import (
	"github.com/golang-jwt/jwt"
	"github.com/irfanguvian/attendance-service/dto"
	"github.com/irfanguvian/attendance-service/entities"
)

type Services struct {
	AuthService AuthService
}

type AuthService interface {
	Login(loginBody dto.LoginBody) (*dto.ResponseLoginService, error)
	Signup(signupBody dto.SignupBody) (string, error)
	SignOut(userID uint) error
	ExchangeToken(refreshToken string) (*dto.ResponseLoginService, error)

	ValidateToken(tokenString string) (jwt.MapClaims, error)

	GetUserByAccessID(accessID string) (*entities.User, error)

	GetAccessTokenByAccessID(accessID string) error
}
