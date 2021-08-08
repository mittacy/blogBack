package apierr

const (
	CodeParamErr = 1		// 前端参数错误
	CodeBackErr = 2			// 后端未知错误
	CodeJsonMarshalErr = 3	// json序列化错误

	// 注册登录
	CodeNoLogin     = 1001
	CodeTokenExpire = 1002

	// 用户
	CodeUserExist   = 2001
	CodeUserNoExist = 2002

	// 分类
	CodeCategoryNameExist = 3001
	CodeCategoryNoExist   = 3002

	// 文章
	CodeArticleNoExist = 4002
)
