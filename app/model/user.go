package model

type User struct {
	Id int64
}

func (*User) TableName() string {
	return "user"
}

