package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/copier"
	"github.com/mittacy/blogBack/apierr"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/app/transform"
	"github.com/mittacy/blogBack/app/validator/articleValidator"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/pkg/response"
	"strconv"
)

type Article struct {
	articleService IArticleService
	transform   transform.Article
	logger      *logger.CustomLogger
}

func NewArticle(articleService IArticleService, logger *logger.CustomLogger) Article {
	return Article{
		articleService: articleService,
		transform: transform.NewArticle(logger),
		logger:    logger,
	}
}

type IArticleService interface {
	Create(article model.Article) (int64, error)
	Delete(id int64) error
	UpdateInfo(article model.Article) error
	UpdateWeight(id int64, weight int64) error
	Get(id int64) (*model.Article, error)
	List(page, pageSize int) ([]model.Article, int64, error)
	ListByCategory(categoryId int64, page, pageSize int) ([]model.Article, int64, error)
	ListHome() ([]model.Article, error)
}

/**
 * @apiVersion 0.1.0
 * @apiGroup Article
 * @api {post} /article 创建文章
 * @apiName Article.Create
 *
 * @apiParam {number} category_id 所属分类id
 * @apiParam {string{1..64}} title 文章标题
 * @apiParam {string{1..1024}} preview_ctx 预览内容
 * @apiParam {string{1..}} content 文章正文
 *
 * @apiSuccess {string} id 创建的文章id
 *
 * @apiSuccessExample {json} Success-Response:
 *     {
 *         "code": 0,
 *         "data": {
 *             "id": 847,
 *         },
 *         "msg": "success"
 *     }
 *
 */
func (ctl *Article) Create(c *gin.Context) {
	req := articleValidator.CreateReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidateErr(c, err)
		return
	}

	article := model.Article{}
	if err := copier.Copy(&article, &req); err != nil {
		response.CopierErrAndLog(c, ctl.logger, err)
		return
	}

	id, err := ctl.articleService.Create(article)
	if err != nil {
		response.CheckErrAndLog(c, ctl.logger, "createArticle", err, apierr.ErrCategoryNoExist)
		return
	}

	res := map[string]int64{
		"id": id,
	}
	response.Success(c, res)
}

/**
 * @apiVersion 0.1.0
 * @apiGroup Article
 * @api {delete} /article/:id 删除文章
 * @apiName Article.Delete
 *
 * @apiParam {number{1..}} id 文章id
 *
 * @apiErrorExample {json} 文章不存在
 *     {
 *       "code": 1,
 *       "msg": "对象不存在",
 *       "data": {}
 *     }
 *
 */
func (ctl *Article) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailMsg(c, "id must be greater than 0")
		return
	}

	if err := ctl.articleService.Delete(id); err != nil {
		response.CheckErrAndLog(c, ctl.logger, "deleteArticle", err, apierr.ErrArticleNoExist, apierr.ErrCategoryNoExist)
		return
	}

	response.Success(c, nil)
}

/**
 * @apiVersion 0.1.0
 * @apiGroup Article
 * @api {patch} /article 编辑文章
 * @apiName Article.Update
 *
 * @apiParam {number=1(文章信息内容),2(文章权重)} update_type 更新类型
 * @apiParam {number{1..}} id 文章id,update_type=1时必须
 * @apiParam {number} category_id 所属分类id,update_type=1时必须
 * @apiParam {string{1..64}} title 文章标题,update_type=1时必须
 * @apiParam {string{1..1024}} preview_ctx 预览内容,update_type=1时必须
 * @apiParam {string{1..}} content 文章正文,update_type=1时必须
 * @apiParam {number{0..}} weight 文章权重,update_type=2时可选
 *
 * @apiErrorExample {json} 文章不存在
 *     {
 *       "code": 1,
 *       "msg": "对象不存在",
 *       "data": {}
 *     }
 *
 */
func (ctl *Article) Update(c *gin.Context) {
	ur := articleValidator.UpdateReq{}
	if err := c.ShouldBindBodyWith(&ur, binding.JSON); err != nil {
		response.ValidateErr(c, err)
		return
	}

	switch ur.UpdateType {
	case 1:
		ctl.updateInfo(c)
	case 2:
		ctl.updateWeight(c)
	}

	return
}

/**
 * @apiVersion 0.1.0
 * @apiGroup Article
 * @api {get} /article/:id 查询文章详情
 * @apiName Article.Get
 *
 * @apiParam {number{1..}} id 文章id
 *
 * @apiSuccess {number} id 文章id
 * @apiSuccess {number} category_id 分类id
 * @apiSuccess {string} category_name 分类名
 * @apiSuccess {string} title 文章标题
 * @apiSuccess {number} views 点击量
 * @apiSuccess {string} created_at 创建时间
 * @apiSuccess {string} updated_at 更新时间
 * @apiSuccess {string} content 文章正文
 *
 * @apiSuccessExample {json} Success-Response:
 *     {
 *         "code": 0,
 *         "data": {
 *             "article": {
 *                 "id": 14,
 *                 "category_id": 5,
 *                 "category_name": "Golang",
 *                 "title": "文章标题",
 *                 "views": 0,
 *                 "created_at": 1625798089,
 *                 "int64": 1625798089,
 *                 "content": "内容",
 *                 "picture": "",
 *                 "sentence": ""
 *             }
 *         },
 *         "msg": "success"
 *     }
 *
 */
func (ctl *Article) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailMsg(c, "id must be greater than 0")
		return
	}

	article, err := ctl.articleService.Get(id)
	if err != nil {
		response.CheckErrAndLog(c, ctl.logger, "getArticle", err, apierr.ErrArticleNoExist)
		return
	}

	ctl.transform.GetReply(c, article)
}

/**
 * @apiVersion 0.1.0
 * @apiGroup Article
 * @api {get} /articles 文章分页列表
 * @apiName Article.List
 *
 * @apiParam {number{1..}} page 页码
 * @apiParam {number{1..50}} page_size 数据分页大小
 * @apiParam {number{1..}} category_id 分类id,传空则为全部
 *
 * @apiSuccess {number} id 文章id
 * @apiSuccess {number} category_id 分类id
 * @apiSuccess {string} category_name 分类名
 * @apiSuccess {string} title 文章标题
 * @apiSuccess {number} views 点击量
 * @apiSuccess {string} preview_ctx 预览内容
 * @apiSuccess {string} created_at 创建时间
 * @apiSuccess {string} updated_at 更新时间
 *
 * @apiSuccessExample {json} Success-Response:
 *     {
 *         "code": 0,
 *         "data": {
 *             "list": [
 *                 {
 *                     "id": 14,
 *                     "category_id": 5,
 *                     "category_name": "Golang",
 *                     "title": "文章标题",
 *                     "views": 0,
 *                     "preview_ctx": "预览内容",
 *                     "created_at": 1625798089,
 *                     "updated_at": 0
 *                 },
 *                 {
 *                     "id": 15,
 *                     "category_id": 7,
 *                     "category_name": "Java",
 *                     "title": "文章标题2",
 *                     "views": 0,
 *                     "preview_ctx": "预览内容2",
 *                     "created_at": 1625798833,
 *                     "updated_at": 0
 *                 }
 *             ],
 *             "total_size": 2
 *         },
 *         "msg": "success"
 *     }
 *
 */
func (ctl *Article) List(c *gin.Context) {
	req := articleValidator.ListReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidateErr(c, err)
		return
	}

	var (
		articles []model.Article
		totalSize int64
		err error
	)

	if req.CategoryId > 0 {
		articles, totalSize, err = ctl.articleService.ListByCategory(req.CategoryId, req.Page, req.PageSize)
	} else {
		articles, totalSize, err = ctl.articleService.List(req.Page, req.PageSize)
	}

	if err != nil {
		response.CheckErrAndLog(c, ctl.logger, "article list", err)
		return
	}

	ctl.transform.ListReply(c, articles, totalSize)
}

/**
 * @apiVersion 0.1.0
 * @apiGroup Article
 * @api {get} /articles_home 首页文章列表
 * @apiName Article.HomeList
 *
 * @apiSuccess {number} id 文章id
 * @apiSuccess {number} category_id 分类id
 * @apiSuccess {string} category_name 分类名
 * @apiSuccess {string} title 文章标题
 * @apiSuccess {number} views 点击量
 * @apiSuccess {string} preview_ctx 预览内容
 * @apiSuccess {string} created_at 创建时间
 * @apiSuccess {string} updated_at 更新时间
 *
 * @apiSuccessExample {json} Success-Response:
 *     {
 *         "code": 0,
 *         "data": {
 *             "list": [
 *                 {
 *                     "id": 14,
 *                     "category_id": 5,
 *                     "category_name": "Golang",
 *                     "title": "文章标题",
 *                     "views": 0,
 *                     "preview_ctx": "预览内容",
 *                     "created_at": 1625798089,
 *                     "updated_at": 0
 *                 },
 *                 {
 *                     "id": 15,
 *                     "category_id": 7,
 *                     "category_name": "Java",
 *                     "title": "文章标题2",
 *                     "views": 0,
 *                     "preview_ctx": "预览内容2",
 *                     "created_at": 1625798833,
 *                     "updated_at": 0
 *                 }
 *             ],
 *             "total_size": 2
 *         },
 *         "msg": "success"
 *     }
 *
 */
func (ctl *Article) HomeList(c *gin.Context) {
	articles, err := ctl.articleService.ListHome()
	if err != nil {
		response.CheckErrAndLog(c, ctl.logger, "home article list", err)
		return
	}

	ctl.transform.HomeListReply(c, articles)
}

func (ctl *Article) updateInfo(c *gin.Context) {
	req := articleValidator.UpdateInfoReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidateErr(c, err)
		return
	}

	article := model.Article{}
	if err := copier.Copy(&article, &req); err != nil {
		response.CopierErrAndLog(c, ctl.logger, err)
		return
	}

	if err := ctl.articleService.UpdateInfo(article); err != nil {
		response.CheckErrAndLog(c, ctl.logger, "update article info", err, apierr.ErrArticleNoExist)
		return
	}

	response.Success(c, nil)
	return
}

func (ctl *Article) updateWeight(c *gin.Context) {
	req := articleValidator.UpdateWeightReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidateErr(c, err)
		return
	}

	article := model.Article{}
	if err := copier.Copy(&article, &req); err != nil {
		response.CopierErrAndLog(c, ctl.logger, err)
		return
	}

	if err := ctl.articleService.UpdateWeight(req.Id, req.Weight); err != nil {
		response.CheckErrAndLog(c, ctl.logger, "update article weight", err, apierr.ErrArticleNoExist)
		return
	}

	response.Success(c, nil)
	return
}

