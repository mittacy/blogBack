package model

type Category struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	ArticleCount int    `json:"article_count"`
}

func (*Category) TableName() string {
	return "category"
}
