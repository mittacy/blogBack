package userValidator

type RegisterReq struct {
	Name      string `json:"name" binding:"required,min=1,max=10"`
	Password  string `json:"password" binding:"required,min=8,max=20"`
	Gender    int8   `json:"gender" binding:"required,oneof=1 5 10"`
	Introduce string `json:"introduce" binding:"omitempty,max=255"`
	Github    string `json:"github" binding:"required,url"`
	Email     string `json:"email" binding:"required,email"`
	Code      string `json:"code" binding:"required,len=6"`
}

type LoginReq struct {
	LoginType int    `json:"login_type" binding:"required,oneof=1 2"`
	Name      string `json:"name" binding:"omitempty,min=1,max=10"`
	Email     string `json:"email" binding:"omitempty,email"`
	Password  string `json:"password" binding:"required,min=8,max=20"`
}

type GetReply struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Gender    int8   `json:"gender"`
	Introduce string `json:"introduce"`
	Github    string `json:"github"`
	Email     string `json:"email"`
	Views     int64  `json:"views"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	LoginAt   int64  `json:"login_at"`
}
