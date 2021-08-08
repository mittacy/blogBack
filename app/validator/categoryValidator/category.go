package categoryValidator

type CreateReq struct {
	Name string `json:"name" binding:"required,min=1,max=16"`
}

type UpdateReq struct {
	UpdateType int `json:"update_type" binding:"required,oneof=1"`
}

type UpdateNameReq struct {
	Id   int64  `json:"id" binding:"required,min=1"`
	Name string `json:"name" binding:"required,min=1,max=16"`
}

type ListReq struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" json:"page_size" binding:"omitempty,min=1"`
}

type ListReply struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	ArticleCount int    `json:"article_count"`
}
