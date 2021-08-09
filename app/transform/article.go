package transform

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/app/validator/articleValidator"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/pkg/response"
)

type Article struct {
	logger *logger.CustomLogger
}

func NewArticle(customLogger *logger.CustomLogger) Article {
	return Article{logger: customLogger}
}

// ArticlePack 数据库数据转化为响应数据
// @param data 数据库数据
// @return reply 响应体数据
// @return err
func (ctl *Article) ArticlePack(data *model.Article) (*articleValidator.GetReply, error) {
	reply := articleValidator.GetReply{}

	if err := copier.Copy(&reply, data); err != nil {
		return nil, err
	}

	return &reply, nil
}

// ArticlesPack 数据库数据转化为响应数据
// @param data 数据库数据
// @return reply 响应体数据
// @return err
func (ctl *Article) ArticlesPack(data []model.Article) (reply []articleValidator.ListReply, err error) {
	err = copier.Copy(&reply, &data)
	return
}

func (ctl *Article) HomeArticlesPack(data []model.Article) (reply []articleValidator.ListHomeReply, err error) {
	err = copier.Copy(&reply, &data)
	return
}

// GetReply 详情响应包装
// @param data 数据库数据
func (ctl *Article) GetReply(c *gin.Context, data *model.Article) {
	reply, err := ctl.ArticlePack(data)
	if err != nil {
		ctl.logger.CopierErrLog(err)
		response.Unknown(c)
		return
	}

	res := map[string]interface{}{
		"article": reply,
	}

	response.Success(c, res)
}

// ListReply 列表响应包装
// @param data 数据库列表数据
// @param totalSize 记录总数(ctl *Article) 
func (ctl *Article) ListReply(c *gin.Context, data []model.Article, totalSize int64) {
	list, err := ctl.ArticlesPack(data)
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

func (ctl *Article) HomeListReply(c *gin.Context, articles []model.Article) {
	list, err := ctl.HomeArticlesPack(articles)
	if err != nil {
		response.CopierErrAndLog(c, ctl.logger, err)
		return
	}

	res := map[string]interface{}{
		"list": list,
	}

	response.Success(c, res)
	return
}

