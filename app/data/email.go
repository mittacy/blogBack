package data

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/mittacy/blogBack/apierr"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/app/service"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/pkg/store/cache"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"strconv"
)

// 实现service层中的data接口

type Email struct {
	db     *gorm.DB
	cache  cache.CustomRedis
	conf   model.EmailConfig
	logger *logger.CustomLogger
}

func NewEmail(db *gorm.DB, cacheConn *redis.Pool, conf model.EmailConfig, logger *logger.CustomLogger) service.IEmailData {
	r := cache.ConnRedisByPool(cacheConn, "email")

	return &Email{
		db:     db,
		cache:  r,
		conf:   conf,
		logger: logger,
	}
}

func (ctl *Email) GetEmailTpl(name string) (*model.EmailTpl, error) {
	tpl := model.EmailTpl{}
	if err := ctl.db.Where("name = ?", name).First(&tpl).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &tpl, nil
}

func (ctl *Email) SaveCode(email string, code string) error {
	const expire = 60 * 5 // 单位:秒

	if _, err := ctl.cache.Do("setex", ctl.cacheCodeKey(email), expire, code); err != nil {
		return err
	}
	return nil
}

func (ctl *Email) GetCode(email string) (string, error) {
	code, err := redis.String(ctl.cache.Do("get", ctl.cacheCodeKey(email)))
	if errors.Is(err, redis.ErrNil) {
		return "", apierr.ErrRegisterCode
	}
	return code, nil
}

func (ctl *Email) InvalidCode(email string) error {
	if err := ctl.cache.Del(ctl.cacheCodeKey(email)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (ctl *Email) SendEmail(mailTo []string, subject string, body string) error {
	port, err := strconv.Atoi(ctl.conf.Port)
	if err != nil {
		return errors.New("端口号必须为数字字符串")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(ctl.conf.User, "Blog"))
	m.SetHeader("To", mailTo...)    // 收件人
	m.SetHeader("Subject", subject) // 邮件主题
	m.SetBody("text/html", body)    // 邮件正文
	d := gomail.NewDialer(ctl.conf.Host, port, ctl.conf.User, ctl.conf.Pass)
	return d.DialAndSend(m)
}

func (ctl *Email) cacheCodeKey(email string) string {
	return fmt.Sprintf("%s:code#%s", ctl.cache.CachePrefixKey(), email)
}
