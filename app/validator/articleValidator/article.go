package articleValidator

type CreateReq struct {
	CategoryId int64  `json:"category_id" binding:"required,min=1"`
	Title      string `json:"title" binding:"required,min=1,max=64"`
	PreviewCtx string `json:"preview_ctx" binding:"required,min=1,max=1024"`
	Content    string `json:"content" binding:"required,min=1"`
}

type UpdateReq struct {
	UpdateType int `json:"update_type" binding:"required,oneof=1 2"`
}

type UpdateInfoReq struct {
	Id         int64  `json:"id" binding:"required,min=1"`
	CategoryId int64  `json:"category_id" binding:"required,min=1"`
	Title      string `json:"title" binding:"required,min=1,max=64"`
	PreviewCtx string `json:"preview_ctx" binding:"required,min=1,max=1024"`
	Content    string `json:"content" binding:"required,min=1"`
}

type UpdateWeightReq struct {
	Id     int64 `json:"id" binding:"required,min=1"`
	Weight int64 `json:"weight" binding:"omitempty,min=0"`
}

type GetReply struct {
	Id           int64  `json:"id"`
	CategoryId   int64  `json:"category_id"`
	CategoryName string `json:"category_name"`
	Title        string `json:"title"`
	Views        int64  `json:"views"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
	Content      string `json:"content"`
	Picture      string `json:"picture"`
	Sentence     string `json:"sentence"`
}

type ListReq struct {
	Page       int   `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize   int   `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=50"`
	CategoryId int64 `form:"category_id" json:"category_id" binding:"omitempty,min=1"`
}

type ListReply struct {
	Id           int64  `json:"id"`
	CategoryId   int64  `json:"category_id"`
	CategoryName string `json:"category_name"`
	Title        string `json:"title"`
	Views        int64  `json:"views"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

type ListHomeReply struct {
	Id           int64  `json:"id"`
	CategoryId   int64  `json:"category_id"`
	CategoryName string `json:"category_name"`
	Title        string `json:"title"`
	Views        int64  `json:"views"`
	PreviewCtx   string `json:"preview_ctx"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}
