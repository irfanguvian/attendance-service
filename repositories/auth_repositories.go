package repositories

import (
	"github.com/irfanguvian/attendance-service/models"
	"gorm.io/gorm"
)

type AuthRepositories struct {
	DB *gorm.DB
}

func NewAuthRepositories(db *gorm.DB) *AuthRepositories {
	return &AuthRepositories{
		DB: db,
	}
}

func (ar *AuthRepositories) CreateUser(email, password string) (uint, error) {
	user := &models.User{
		Email:    email,
		Password: password,
	}
	if err := ar.DB.Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (ar *AuthRepositories) CreateAccessToken(userID uint, accessID string) error {
	accessTokenModel := &models.AccessToken{
		UserID: userID,
		ID:     accessID,
	}
	return ar.DB.Create(accessTokenModel).Error
}

func (ar *AuthRepositories) CreateRefreshToken(accessID string, refreshToken string) error {
	refreshTokenModel := &models.RefreshToken{
		AccessTokenID: accessID,
		ID:            refreshToken,
	}
	return ar.DB.Create(refreshTokenModel).Error
}

func (ar *AuthRepositories) DeleteAccessTokenByUserID(userID uint) error {
	accessToken := &models.AccessToken{UserID: userID}
	return ar.DB.Where(accessToken).Delete(&models.AccessToken{}).Error
}

func (ar *AuthRepositories) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := ar.DB.Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ar *AuthRepositories) GetRefreshTokenByAccessID(accessID string) (string, error) {
	refreshToken := &models.RefreshToken{}
	if err := ar.DB.Where("access_token_id = ?", accessID).First(refreshToken).Error; err != nil {
		return "", err
	}
	return refreshToken.ID, nil
}

func (ar *AuthRepositories) GetUserByAccessID(accessID string) (*models.User, error) {
	accessToken := &models.AccessToken{}
	if err := ar.DB.Where("id = ?", accessID).First(accessToken).Error; err != nil {
		return nil, err
	}

	user := &models.User{}
	if err := ar.DB.Where("id = ?", accessToken.UserID).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
