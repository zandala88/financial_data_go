package middleware

import (
	"financia/util"
	"github.com/gin-gonic/gin"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		userClaims, err := util.VerifyJWT(token)
		if err != nil {
			c.Abort()
			util.FailRespWithCode(c, util.InvalidToken)
			return
		}
		c.Set("user_id", userClaims.UserId)
		c.Next()
	}
}
