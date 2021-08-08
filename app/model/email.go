package model

const EmailRegisterTplName = "register_code"

type EmailTpl struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type EmailConfig struct {
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func (*EmailTpl) TableName() string {
	return "email_tpl"
}

