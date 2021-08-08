package model

type Article struct {
	Id           int64  `json:"id"`
	Weight       int64  `json:"weight"`
	CategoryId   int64  `json:"category_id"`
	CategoryName string `json:"category_name" gorm:"-"`
	Title        string `json:"title"`
	Views        int64  `json:"views"`
	PreviewCtx   string `json:"preview_ctx"`
	CreatedAt    int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    int64  `json:"updated_at" gorm:"autoUpdateTime"`
	Content      string `json:"content"`
	Deleted      int8   `json:"deleted"`
	Picture      string `json:"picture"`
	Sentence     string `json:"sentence"`
}

func (*Article) TableName() string {
	return "article"
}

const (
	ArticleDeletedNo  = 0
	ArticleDeletedYes = 1
)

const (
	_ = iota * 5
	WeightLow
	WeightSecond
	WeightThree
	WeightFour
	WeightFive
)

