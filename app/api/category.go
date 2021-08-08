package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/copier"
	"github.com/mittacy/blogBack/apierr"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/app/transform"
	"github.com/mittacy/blogBack/app/validator/categoryValidator"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/pkg/response"
	"strconv"
)

/**
 * @apiDefine CategoryNameExist 分类名已存在
 * @apiErrorExample {json} 分类名已存在
 *     {
 *       "code": 1,
 *       "msg": "分类名已存在",
 *       "data": {}
 *     }
 */

type Category struct {
	categoryService ICategoryService
	transform   transform.Category
	logger      *logger.CustomLogger
}

func NewCategory(categoryService ICategoryService, logger *logger.CustomLogger) Category {
	return Category{
		categoryService: categoryService,
		transform: transform.NewCategory(logger),
		logger:    logger,
	}
}

type ICategoryService interface {
	Create(category model.Category) (int64, error)
	Delete(id int64) error
	UpdateName(category model.Category) error
	List(page, pageSize int) ([]model.Category, int, error)
}

/**
 * @apiVersion 0.1.0
 * @apiGroup Category
 * @api {post} /category 创建文章分类
 * @apiName Category.Create
 *
 * @apiParam {string{1..16}} name 分类名
 *
 * @apiSuccess {string} id 创建的分类id
 *
 * @apiSuccessExample {json} Success-Response:
 *     {
 *         "code": 0,
 *         "data": {
 *             "id": 2,
 *         },
 *         "msg": "success"
 *     }
 *
 * @apiUse CategoryNameExist
 *
 */
func (ctl *Category) Create(c *gin.Context) {
	req := categoryValidator.CreateReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidateErr(c, err)
		return
	}

	category := model.Category{}
	if err := copier.Copy(&category, &req); err != nil {
		response.CopierErrAndLog(c, ctl.logger, err)
		return
	}

	categoryId, err := ctl.categoryService.Create(category)
	if err != nil {
		response.CheckErrAndLog(c, ctl.logger, "createCategory", err, apierr.ErrCategoryNameExist)
		return
	}

	res := map[string]interface{}{
		"id": categoryId,
	}

	response.Success(c, res)
	return
}

/**
 * @apiVersion 0.1.0
 * @apiGroup Category
 * @api {delete} /category/:id 删除文章分类
 * @apiName Category.Delete
 *
 * @apiParam {number} id 分类id
 *
 * @apiErrorExample {json} 分类不存在
 *     {
 *       "code": 1,
 *       "msg": "对象不存在",
 *       "data": {}
 *     }
 *
 */
func (ctl *Category) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.FailMsg(c, "id必须大于0")
		return
	}

	if err := ctl.categoryService.Delete(id); err != nil {
		response.CheckErrAndLog(c, ctl.logger, "deleteCategory", err, apierr.ErrCategoryNoExist)
		return
	}

	response.Success(c, nil)
	return
}

/**
 * @apiVersion 0.1.0
 * @apiGroup Category
 * @api {patch} /category 更新分类
 * @apiName Category.Update
 *
 * @apiParam {number=1(分类名)} update_type 更新类型
 * @apiParam {number} id 分类id
 * @apiParam {string} string 分类新名字,update_type=1时必须
 *
 * @apiErrorExample {json} 分类不存在
 *     {
 *       "code": 1,
 *       "msg": "对象不存在",
 *       "data": {}
 *     }
 *
 * @apiUse CategoryNameExist
 *
 */
func (ctl *Category) Update(c *gin.Context) {
	req := categoryValidator.UpdateReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidateErr(c, err)
		return
	}

	switch req.UpdateType {
	case 1:
		ctl.updateName(c)
	}

	return
}

/**
 * @apiVersion 0.1.0
 * @apiGroup Category
 * @api {get} /categories 获取分类列表
 * @apiName Category.List
 *
 * @apiParam {number{1..}} page 页码
 * @apiParam {number{0..}} page_size 请求的数据分页大小,0表示不分页获取全部
 *
 * @apiSuccessExample {json} Success-Response:
 *     {
 *         "code": 0,
 *         "data": {
 *             "total_size": 2,
 *             "list": [
 *                 {
 *                     "id": 2,
 *                     "name": "Java",
 *                     "article_count": 0
 *                 },
 *                 {
 *                     "id": 3,
 *                     "name": "Vue",
 *                     "article_count": 0
 *                 }
 *             ]
 *         },
 *         "msg": "success"
 *     }
 */
func (ctl *Category) List(c *gin.Context) {
	req := categoryValidator.ListReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidateErr(c, err)
		return
	}

	categories, count, err := ctl.categoryService.List(req.Page, req.PageSize)
	if err != nil {
		response.CheckErrAndLog(c, ctl.logger, "categoriesList", err)
		return
	}

	ctl.transform.ListReply(c, categories, int64(count))
}

func (ctl *Category) updateName(c *gin.Context) {
	req := categoryValidator.UpdateNameReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidateErr(c, err)
		return
	}

	category := model.Category{}
	if err := copier.Copy(&category, &req); err != nil {
		response.CopierErrAndLog(c, ctl.logger, err)
		return
	}

	if err := ctl.categoryService.UpdateName(category); err != nil {
		response.CheckErrAndLog(c, ctl.logger, "updateCategoryName", err, apierr.ErrCategoryNoExist, apierr.ErrCategoryNameExist)
		return
	}

	response.Success(c, nil)
}

