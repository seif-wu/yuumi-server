package user

import (
	"yuumi/api/auth"

	"github.com/gin-gonic/gin"
)

func (u *handler) Current(c *gin.Context) {
	user, _ := c.Get(auth.IdentityKey)

	c.JSON(200, gin.H{
		"success": true,
		"data":    user,
	})
}
