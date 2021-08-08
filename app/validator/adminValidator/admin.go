package adminValidator

type AdminLoginReq struct {
	Name     string `json:"name" binding:"required,min=1,max=10"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

