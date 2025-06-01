package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/irfanguvian/attendance-service/dto"
	"github.com/irfanguvian/attendance-service/entities"
	"github.com/irfanguvian/attendance-service/interfaces"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	Repositories interfaces.Repositories
}

func NewAuthService(repo interfaces.Repositories) interfaces.AuthService {
	return &AuthService{
		Repositories: repo,
	}
}

func (as *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}

func CreateToken(args jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, args)

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (as *AuthService) Signup(signupBody dto.SignupBody) (string, error) {
	// check email first
	checkEmail, err := as.Repositories.UserRepository.GetUserByEmail(signupBody.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return "server error", err // Handle error appropriately
	}

	if checkEmail != nil {
		return "email already exists", errors.New("email already exists") // Handle error appropriately
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return "failed to hash password", err // Handle error appropriately
	}

	_, err = as.Repositories.UserRepository.CreateUser(signupBody.Email, string(hashedPassword))
	if err != nil {
		return "failed to create user", err // Handle error appropriately
	}
	return "Success", nil // Return success message or user ID as needed
}

func (as *AuthService) Login(loginBody dto.LoginBody) (*dto.ResponseLoginService, error) {
	fmt.Println("Login attempt with email:", loginBody.Email)
	getUser, err := as.Repositories.UserRepository.GetUserByEmail(loginBody.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("email not found")
		}
		return nil, errors.New("server error")
	}

	accessTokenID := uuid.New().String()
	refreshTokenID := uuid.New().String()

	err = bcrypt.CompareHashAndPassword([]byte(getUser.Password), []byte(loginBody.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	// delete token existing
	err = as.Repositories.UserRepository.DeleteAccessTokenByUserID(getUser.ID)
	if err != nil {
		return nil, errors.New("failed to delete existing access token")
	}

	claims := jwt.MapClaims{
		"access_id": accessTokenID,
		"exp":       jwt.TimeFunc().Add(1 * time.Minute).Unix(), // Set expiration time to 24 hours
	}

	claimsRefreshToken := jwt.MapClaims{
		"access_id":  accessTokenID,
		"refresh_id": refreshTokenID,
		"exp":        jwt.TimeFunc().Add(10 * time.Hour).Unix(), // Set expiration time to 24 hours
	}

	token, err := CreateToken(claims)
	if err != nil {
		return nil, err
	}

	tokenRefresh, err := CreateToken(claimsRefreshToken)
	if err != nil {
		return nil, err
	}

	err = as.Repositories.UserRepository.CreateAccessToken(getUser.ID, accessTokenID)

	if err != nil {
		return nil, err
	}

	err = as.Repositories.UserRepository.CreateRefreshToken(accessTokenID, refreshTokenID)
	if err != nil {
		return nil, err
	}

	return &dto.ResponseLoginService{
		AccessToken:  token,
		RefreshToken: tokenRefresh,
	}, nil
}

func (as *AuthService) SignOut(userID uint) error {
	err := as.Repositories.UserRepository.DeleteAccessTokenByUserID(userID)
	if err != nil {
		return errors.New("failed to delete access token")
	}
	return nil
}

func (as *AuthService) ExchangeToken(refreshToken string) (*dto.ResponseLoginService, error) {
	claims, err := as.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	_, err = as.Repositories.UserRepository.GetRefreshTokenByAccessID(claims["access_id"].(string))
	if err != nil {
		return nil, errors.New("failed to get refresh token")
	}

	user, err := as.Repositories.UserRepository.GetUserByAccessID(claims["access_id"].(string))
	if err != nil {
		return nil, errors.New("failed extracting user from access ID")
	}

	err = as.Repositories.UserRepository.DeleteAccessTokenByUserID(user.ID)
	if err != nil {
		return nil, errors.New("failed to delete existing access token")
	}

	newAccessID := uuid.New().String()
	newRefreshID := uuid.New().String()

	newClaims := jwt.MapClaims{
		"access_id": newAccessID,
	}

	newClaimsRefresh := jwt.MapClaims{
		"access_id":  newAccessID,
		"refresh_id": newRefreshID,
	}

	newAccessToken, err := CreateToken(newClaims)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := CreateToken(newClaimsRefresh)
	if err != nil {
		return nil, err
	}

	err = as.Repositories.UserRepository.CreateAccessToken(user.ID, newAccessID)
	if err != nil {
		return nil, err
	}

	err = as.Repositories.UserRepository.CreateRefreshToken(newAccessID, newRefreshID)
	if err != nil {
		return nil, err
	}

	return &dto.ResponseLoginService{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (as *AuthService) GetUserByAccessID(accessID string) (*entities.User, error) {
	user, err := as.Repositories.UserRepository.GetUserByAccessID(accessID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("server error")
	}
	return &entities.User{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

func (as *AuthService) GetAccessTokenByAccessID(accessID string) error {
	_, err := as.Repositories.UserRepository.GetRefreshTokenByAccessID(accessID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("access token not found")
		}
		return errors.New("server error")
	}
	return nil
}
