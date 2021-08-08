package model

type Admin struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

func (*Admin) TableName() string {
	return "admin"
}

