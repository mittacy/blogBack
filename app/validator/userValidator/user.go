package userValidator

type CreateReq struct {}
type CreateReply struct{}

type DeleteReq struct {}
type DeleteReply struct{}

type UpdateReq struct {
	UpdateType int `json:"update_type" binding:"required,oneof=1"`
}

type UpdateInfoReq struct {}
type UpdateInfoReply struct{}

type GetReq struct {}
type GetReply struct {}

type ListReq struct {}
type ListReply struct {}

