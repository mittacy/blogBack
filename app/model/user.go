package model

type User struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Salt      string `json:"salt"`
	Gender    int8   `json:"gender"`
	Introduce string `json:"introduce"`
	Github    string `json:"github"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
	LoginAt   int64  `json:"login_at" gorm:"autoCreateTime"`
}

func (*User) TableName() string {
	return "user"
}

const (
	// 性别
	UserGenderSecret = 1  // 性别-保密
	UserGenderBoy    = 5  // 性别-男
	UserGenderGirl   = 10 // 性别-女

	// 唯一索引名字
	UserIdxName  = "uidx_name"  // name索引名
	UserIdxEmail = "uidx_email" // email索引名

	// 登录方式
	LoginTypeByName  = 1 // 使用用户名登录
	LoginTypeByEmail = 2 // 使用邮箱登录

	// 用户身份
	UserRoleNormal = 1  // 普通成员
	UserRoleAdmin  = 10 // 管理员
)

