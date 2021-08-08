package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/pkg/response"
)

const (
	ActionAddCategory = iota
	ActionPutCategory
	ActionDeleteCategory

	ActionAddArticle
	ActionPutArticle
	ActionDeleteArticle
)

func Operate(action int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetInt("role")

		switch action {
		case ActionAddCategory, ActionPutCategory, ActionDeleteCategory, ActionAddArticle, ActionPutArticle, ActionDeleteArticle:
			if userRole < model.UserRoleAdmin {
				response.FailMsg(c, "权限不足")
				c.Abort()
				return
			}
		default:
			response.FailMsg(c, "操作参数有误")
			c.Abort()
			return
		}

		c.Next()
	}
}
