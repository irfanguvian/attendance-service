package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/irfanguvian/attendance-service/dto"
)

func GetContext(c *gin.Context) *dto.ContextUser {
	var ctx dto.ContextUser
	if c == nil {
		return nil
	}
	user_id, isExist := c.Get("user_id")
	if !isExist {
		return nil
	}
	email, isExist := c.Get("email")
	if !isExist {
		return nil
	}
	ctx.UserID = uint(user_id.(uint))
	ctx.Email = email.(string)
	return &ctx
}
