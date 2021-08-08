package emailValidator

type RegisterCodeReq struct {
	Email string `form:"email" binding:"required,email"`
}

