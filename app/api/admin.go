package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mittacy/blogBack/apierr"
	"github.com/mittacy/blogBack/app/transform"
	"github.com/mittacy/blogBack/app/validator/adminValidator"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/pkg/response"
)

type Admin struct {
	adminService IAdminService
	transform   transform.Admin
	logger      *logger.CustomLogger
}

func NewAdmin(adminService IAdminService, logger *logger.CustomLogger) Admin {
	return Admin{
		adminService: adminService,
		transform: transform.NewAdmin(logger),
		logger:    logger,
	}
}

type IAdminService interface {
	Login(adminValidator.AdminLoginReq) (string, error)
}

/**
 * @apiVersion 0.1.0
 * @apiGroup Admin
 * @api {post} /session/admin/login 管理员登录
 * @apiName Admin.Login
 *
 * @apiParam {string{1..10}} name 用户名
 * @apiParam {string{2..20}} password 密码
 *
 * @apiSuccess {string} token 登录身份token
 *
 * @apiSuccessExample {json} Success-Response:
 *     {
 *         "code": 0,
 *         "data": {
 *           "token": "xxx"
 *         },
 *         "msg": "success"
 *     }
 *
 * @apiErrorExample {json} 账号或密码错误:
 *     {
 *       "code": 1,
 *       "msg": "账号或密码错误",
 *       "data": {}
 *     }
 */
func (ctl *Admin) Login(c *gin.Context) {
	var req adminValidator.AdminLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidateErr(c, err)
		return
	}

	token, err := ctl.adminService.Login(req)
	if err != nil {
		response.CheckErrAndLog(c, ctl.logger, "adminLogin", err, apierr.ErrUserOrPassword)
		return
	}

	response.Success(c, map[string]string{"token": token})
}

