package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/irfanguvian/attendance-service/dto"
	"github.com/irfanguvian/attendance-service/interfaces"
	"github.com/irfanguvian/attendance-service/utils"
)

type AuthController struct {
	AuthService interfaces.AuthService
}

func NewAuthController(authService interfaces.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

func (ac *AuthController) Signup(c *gin.Context) {
	// Call the AuthService's Signup method
	var reqBody dto.SignupBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	errMsg, err := ac.AuthService.Signup(reqBody)
	if err != nil {
		utils.ErrorResponse(c, 500, errMsg)
		return
	}
	utils.SuccessResponse(c, 200, "Signup successful", nil)
}

func (ac *AuthController) Login(c *gin.Context) {
	// Call the AuthService's Login method
	var reqBody dto.LoginBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	token, err := ac.AuthService.Login(reqBody)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, 200, "Login successful", token)
}

func (ac *AuthController) SignOut(c *gin.Context) {
	context := utils.GetContext(c)
	fmt.Print("SignOut called with user ID: ", context.UserID)
	err := ac.AuthService.SignOut(context.UserID)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, 200, "Logout successful", nil)
}

func (ac *AuthController) ExchangeToken(c *gin.Context) {
	var reqBody dto.ExchangeTokenBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	newToken, err := ac.AuthService.ExchangeToken(reqBody.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, 200, "Token exchanged successfully", newToken)
}
