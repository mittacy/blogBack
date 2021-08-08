package transform

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/app/validator/categoryValidator"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/pkg/response"
)

type Category struct {
	logger *logger.CustomLogger
}

func NewCategory(customLogger *logger.CustomLogger) Category {
	return Category{logger: customLogger}
}

// CategorysPack 数据库数据转化为响应数据
// @param data 数据库数据
// @return reply 响应体数据
// @return err
func (ctl *Category) CategoriesPack(data []model.Category) (reply []categoryValidator.ListReply, err error) {
	err = copier.Copy(&reply, &data)
	return
}

// ListReply 列表响应包装
// @param data 数据库列表数据
// @param totalSize 记录总数(ctl *Category) 
func (ctl *Category) ListReply(c *gin.Context, data []model.Category, totalSize int64) {
	list, err := ctl.CategoriesPack(data)
	if err != nil {
		response.CopierErrAndLog(c, ctl.logger, err)
		return
	}

	res := map[string]interface{}{
		"list":       list,
		"total_size": totalSize,
	}

	response.Success(c, res)
}

