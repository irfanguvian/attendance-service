package dto

type LoginBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ContextUser struct {
	UserID uint     `json:"user_id"`
	Email  string   `json:"email"`
}

type SignupBody struct {
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required"`
}

type ResponseLoginService struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`	
}

type ExchangeTokenBody struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
