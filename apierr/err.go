package apierr

import (
	"errors"
)

var (
	ErrParam       = errors.New("参数错误")
	ErrCopier      = errors.New("结构体转化错误")
	ErrJsonMarshal = errors.New("json序列化错误")

	// 缓存
	ErrCacheNoExist = errors.New("查询的缓存不存在")

	// 注册登录
	ErrUserEmailExist = errors.New("邮箱已注册")
	ErrRegisterCode   = errors.New("验证码不正确")
	ErrUserNameExist  = errors.New("name已被占用")
	ErrUserOrPassword = errors.New("账号或密码错误")
	ErrLoginExpire    = errors.New("登录信息过期")
	ErrNoLogin        = errors.New("未登录")

	// 用户
	ErrUserNoExist = errors.New("用户不存在")

	// 分类
	ErrCategoryNameExist = errors.New("分类名已存在")
	ErrCategoryNoExist   = errors.New("分类不存在")

	// 文章
	ErrArticleNoExist = errors.New("文章不存在")
)

var errCode = map[error]int{
	ErrParam:  CodeParamErr,
	ErrCopier: CodeBackErr,
	ErrJsonMarshal: CodeJsonMarshalErr,

	// 注册登录
	ErrLoginExpire: CodeTokenExpire,
	ErrNoLogin:     CodeNoLogin,

	// 用户
	ErrUserNoExist: CodeUserNoExist,

	// 分类
	ErrCategoryNameExist: CodeCategoryNameExist,
	ErrCategoryNoExist:   CodeCategoryNoExist,

	// 文章
	ErrArticleNoExist: CodeArticleNoExist,
}

func ErrCode(err error) int {
	if v, ok := errCode[err]; ok {
		return v
	}
	return CodeParamErr
}
