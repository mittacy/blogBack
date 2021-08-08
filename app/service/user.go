package service

import (
	"github.com/mittacy/blogBack/apierr"
	"github.com/mittacy/blogBack/app/api"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/pkg/jwt"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/utils"
	"github.com/pkg/errors"
	"time"
)

type User struct {
	userData  IUserData
	emailData IEmailData
	logger    *logger.CustomLogger
}

// 编写实现api层中的各个service接口的构建方法

func NewUser(userData IUserData, emailData IEmailData, logger *logger.CustomLogger) api.IUserService {
	return &User{
		userData:  userData,
		emailData: emailData,
		logger:    logger,
	}
}

type IUserData interface {
	Create(*model.User) error
	Get(id int64) (*model.User, error)
	GetByName(name string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	UpdatesById(user model.User, updateFields []string, isCleanCache bool) error
}

func (ctl *User) Register(user model.User, code string) (userId int64, err error) {
	// 1. 检查验证码
	rCode, err := ctl.emailData.GetCode(user.Email)
	if err != nil {
		return
	}

	if rCode != code {
		err = apierr.ErrRegisterCode
		return
	}

	// 2. 加密密码
	user.Password, user.Salt = utils.Encryption(user.Password)

	// 3. 创建用户
	if err = ctl.userData.Create(&user); err != nil {
		return
	}

	// 4. 让验证码失效
	if err := ctl.emailData.InvalidCode(user.Email); err != nil {
		ctl.logger.Sugar().Errorf("置位注册验证码失效错误: %s", err)
	}

	return user.Id, nil
}

func (ctl *User) GetUserInfo(id int64) (*model.User, error) {
	return ctl.userData.Get(id)
}

func (ctl *User) LoginByName(name, password string) (string, error) {
	user := model.User{Name: name, Password: password}
	return ctl.login(model.LoginTypeByName, user)
}

func (ctl *User) LoginByEmail(email, password string) (string, error) {
	user := model.User{Email: email, Password: password}
	return ctl.login(model.LoginTypeByEmail, user)
}

func (ctl *User) login(loginType int, user model.User) (token string, err error) {
	// 1. 查询数据库中用户的信息
	realUser := &model.User{}

	switch loginType {
	case model.LoginTypeByName:
		realUser, err = ctl.userData.GetByName(user.Name)
	case model.LoginTypeByEmail:
		realUser, err = ctl.userData.GetByEmail(user.Email)
	}

	if err != nil {
		if errors.Is(err, apierr.ErrUserNoExist) { // 隐藏错误信息，不让用户知道是账号不存在
			err = apierr.ErrUserOrPassword
		}
		return
	}

	// 2. 校验密码
	if utils.EncryptionBySalt(user.Password, realUser.Salt) != realUser.Password {
		return
	}

	// 3. 更新登录时间
	u := model.User{Id: realUser.Id, LoginAt: time.Now().Unix()}
	go func() {
		if err := ctl.userData.UpdatesById(u, []string{"login_at"}, false); err != nil {
			ctl.logger.Sugar().Errorf("update login_at err: %s", err)
		}
	}()

	// 4. 生成 token
	token, err = jwt.Token.Create(user.Id, model.UserRoleNormal)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return token, nil
}
