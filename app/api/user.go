package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/mittacy/blogBack/apierr"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/app/transform"
	"github.com/mittacy/blogBack/app/validator/userValidator"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/pkg/response"
	"strconv"
)

type User struct {
	userService IUserService
	transform   transform.User
	logger      *logger.CustomLogger
}

func NewUser(userService IUserService, logger *logger.CustomLogger) User {
	return User{
		userService: userService,
		transform: transform.NewUser(logger),
		logger:    logger,
	}
}

type IUserService interface {
	Register(user model.User, code string) (int64, error)
	GetUserInfo(id int64) (*model.User, error)
	LoginByName(name, password string) (string, error)
	LoginByEmail(email, password string) (string, error)
}

/**
 * @apiVersion 0.1.0
 * @apiGroup User
 * @api {post} /user 用户注册
 * @apiName User.Register
 *
 * @apiParam {string{1..10}} name 用户名
 * @apiParam {string{2..20}} password 用户密码
 * @apiParam {number=1,5,10} gender 用户性别(1保密,5男,10女)
 * @apiParam {string{..255}} introduce="空" 个人介绍
 * @apiParam {string} github github地址
 * @apiParam {string} email 邮箱
 * @apiParam {string} code 邮箱验证码
 *
 * @apiSuccess {string} id 创建的用户id
 * @apiSuccessExample {json} Success-Response:
 *     {
 *         "code": 0,
 *         "data": {
 *             "id": 237,
 *         },
 *         "msg": "success"
 *     }
 *
 * @apiErrorExample {json} 用户名被占用
 *     {
 *       "code": 1,
 *       "msg": "用户名被占用",
 *       "data": {}
 *     }
 * @apiErrorExample {json} 邮箱被注册
 *     {
 *       "code": 1,
 *       "msg": "邮箱被注册",
 *       "data": {}
 *     }
 * @apiErrorExample {json} 邮箱验证码错误
 *     {
 *       "code": 1,
 *       "msg": "验证码错误",
 *       "data": {}
 *     }
 *
 */
func (ctl *User) Register(c *gin.Context) {
	register := userValidator.RegisterReq{}
	if err := c.ShouldBindJSON(&register); err != nil {
		response.ValidateErr(c, err)
		return
	}

	user := model.User{}
	if err := copier.Copy(&user, &register); err != nil {
		response.CopierErrAndLog(c, ctl.logger, err)
		return
	}

	userId, err := ctl.userService.Register(user, register.Code)
	if err != nil {
		response.CheckErrAndLog(c, ctl.logger, "registerUser", err, apierr.ErrUserNameExist, apierr.ErrUserEmailExist, apierr.ErrRegisterCode)
		return
	}

	response.Success(c, map[string]int64{"id": userId})
}

/**
 * @apiVersion 0.1.0
 * @apiGroup User
 * @api {post} /session/user/login 用户登录
 * @apiName User.Login
 *
 * @apiParam {number} login_type 登录方式(1:昵称 2:邮箱)
 * @apiParam {string{1..10}} name 用户名,login_type=1时必须
 * @apiParam {string} email 邮箱,login_type=2时必须
 * @apiParam {string{2..20}} password 用户密码
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
 *
 */
func (ctl *User) Login(c *gin.Context) {
	userLogin := userValidator.LoginReq{}
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		response.ValidateErr(c, err)
		return
	}

	var token string
	var err error

	switch userLogin.LoginType {
	case 1:
		token, err = ctl.userService.LoginByName(userLogin.Name, userLogin.Password)
	case 2:
		token, err = ctl.userService.LoginByEmail(userLogin.Email, userLogin.Password)
	default:
		response.FailMsg(c, "update_type param err")
		return
	}

	if err != nil {
		response.CheckErrAndLog(c, ctl.logger, "userLogin", err, apierr.ErrUserOrPassword)
		return
	}

	response.Success(c, map[string]string{"token": token})
}

/**
 * @apiVersion 0.1.0
 * @apiGroup User
 * @api {get} /user/{id} 获取用户信息
 * @apiName User.GetInfo
 *
 * @apiSuccess {number} id 用户id
 * @apiSuccess {string} name 用户昵称
 * @apiSuccess {string} gender 用户性别(1保密,5男,10女)
 * @apiSuccess {string} introduce 个人简介
 * @apiSuccess {string} github github地址
 * @apiSuccess {string} email 邮箱地址
 * @apiSuccess {string} created_at 创建时间
 * @apiSuccess {string} login_at 最后登录时间
 *
 * @apiSuccessExample {json} Success-Response:
 * {
 *     "code": 0,
 *     "data": {
 *         "user": {
 *             "id": 14,
 *             "name": "mittacy",
 *             "gender": 5,
 *             "introduce": "介绍",
 *             "github": "https://www.github.com",
 *             "email": "mail@mittacy.com",
 *             "created_at": "2021-06-30T03:50:11Z",
 *             "login_at": "2021-06-30T03:50:11Z"
 *         },
 *     },
 *     "msg": "success"
 * }
 *
 *
 * @apiErrorExample {json} 用户不存在
 *     {
 *       "code": 1002,
 *       "msg": "用户不存在",
 *       "data": {}
 *     }
 *
 */
func (ctl *User) GetInfo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidateErr(c, err)
		return
	}

	user, err := ctl.userService.GetUserInfo(id)
	if err != nil {
		response.CheckErrAndLog(c, ctl.logger, "getUser", err, apierr.ErrUserNoExist)
		return
	}

	ctl.transform.GetReply(c, user)
}
