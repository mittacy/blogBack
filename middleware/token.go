package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mittacy/blogBack/apierr"
	"github.com/mittacy/blogBack/pkg/jwt"
	"github.com/mittacy/blogBack/pkg/response"
)

// ParseToken 解析token到gin.Context
func ParseToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			response.FailErr(c, apierr.ErrNoLogin)
			c.Abort()
			return
		}

		token, err := jwt.Token.Parse(accessToken)
		if token == nil || err != nil {
			response.FailErr(c, apierr.ErrLoginExpire)
			c.Abort()
			return
		}

		c.Set("userId", token.UserId)
		c.Set("role", token.Role)

		c.Next()
	}
}
