package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/irfanguvian/attendance-service/interfaces"
	"github.com/irfanguvian/attendance-service/utils"
)

type AuthMiddleware struct {
	AuthService interfaces.AuthService
}

func NewAuthMiddleware(authService interfaces.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService: authService,
	}
}

func (up *AuthMiddleware) ProtecHandlerRequest(context *gin.Context) {
	tokenString := context.Request.Header.Get("Authorization")
	if tokenString == "" {
		context.JSON(401, gin.H{
			"success": false,
			"message": "Token is required",
		})
		context.Abort()
		return
	}

	claims, err := up.AuthService.ValidateToken(tokenString)
	if err != nil {
		context.JSON(401, gin.H{
			"success": false,
			"message": "Invalid Token",
		})
		context.Abort()
		return
	}
	err = up.AuthService.GetAccessTokenByAccessID(claims["access_id"].(string))
	if err != nil {
		context.JSON(401, gin.H{
			"success": false,
			"message": "Access token not found",
		})
		context.Abort()
		return
	}
	// get user info
	getUser, err := up.AuthService.GetUserByAccessID(claims["access_id"].(string))
	if err != nil {
		utils.ErrorResponse(context, 500, "user not found")
		context.Abort()
		return
	}
	context.Set("email", getUser.Email)
	context.Set("user_id", getUser.ID)
	context.Next()
}
