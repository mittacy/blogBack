package service

import (
	"github.com/mittacy/blogBack/app/api"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/pkg/logger"
	"go.uber.org/zap"
)

type Article struct {
	articleData  IArticleData
	categoryData IArticleCategoryData
	logger       *logger.CustomLogger
}

// 编写实现api层中的各个service接口的构建方法

func NewArticle(articleData IArticleData, categoryData IArticleCategoryData, logger *logger.CustomLogger) api.IArticleService {
	return &Article{
		articleData: articleData,
		categoryData: categoryData,
		logger:      logger,
	}
}

type IArticleData interface {
	Insert(article *model.Article) error
	Delete(id int64) error
	UpdateById(article *model.Article, updateFields []string) error
	Get(id int64) (*model.Article, error)
	GetSum() (int64, error)
	GetSumByCategory(categoryId int64) (int64, error)
	List(selectFields []string, page, pageSize int) ([]model.Article, error)
	ListByCategory(selectFields []string, categoryId int64, page, pageSize int) ([]model.Article, error)
	ListByWeight(selectFields []string, count int) ([]model.Article, error)
	IncrView(id int64) error
}

type IArticleCategoryData interface {
	ExpireCategoryData()
	GetCategoriesMap() (map[int64]model.Category, error)
}

func (ctl *Article) Create(article model.Article) (int64, error) {
	if err := ctl.articleData.Insert(&article); err != nil {
		return 0, err
	}

	// 让全部分类缓存失效
	ctl.categoryData.ExpireCategoryData()

	return article.Id, nil
}

func (ctl *Article) Delete(id int64) error {
	if err := ctl.articleData.Delete(id); err != nil {
		return err
	}

	// 让全部分类缓存失效
	ctl.categoryData.ExpireCategoryData()

	return nil
}

func (ctl *Article) UpdateInfo(article model.Article) error {
	fields := []string{"category_id", "title", "preview_ctx", "content"}
	if err := ctl.articleData.UpdateById(&article, fields); err != nil {
		return err
	}
	return nil
}

func (ctl *Article) UpdateWeight(id, weight int64) error {
	article := model.Article{Id: id, Weight: weight}
	fields := []string{"weight"}

	if err := ctl.articleData.UpdateById(&article, fields); err != nil {
		return err
	}

	return nil
}

func (ctl *Article) Get(id int64) (*model.Article, error) {
	/*
	 * 1. 获取文章
	 * 2. 查询文章所属分类的分类名字
	 * 3. 文章阅读量+1
	 */
	article, err := ctl.articleData.Get(id)
	if err != nil {
		return nil, err
	}

	// 填充文章的分类名
	categories, err := ctl.categoryData.GetCategoriesMap()
	if err != nil {
		return nil, err
	}
	article.CategoryName = categories[article.CategoryId].Name

	// 文章阅读量+1
	if err := ctl.articleData.IncrView(id); err != nil {
		zap.S().Errorf("article incr view, err: %s", err)
	}

	return article, nil
}

func (ctl *Article) List(page, pageSize int) ([]model.Article, int64, error) {
	/*
	 * 1. 获取文章列表
	 * 2. 填充文章的分类信息
	 * 3. 查询文章总记录量
	 */
	fields := []string{"id", "category_id", "title", "views", "created_at", "updated_at"}

	articles, err := ctl.articleData.List(fields, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	if err := ctl.FillArticlesCategoryName(articles); err != nil {
		return nil, 0, err
	}

	totalSize, err := ctl.articleData.GetSum()
	if err != nil {
		return nil, 0, err
	}

	return articles, totalSize, nil
}

func (ctl *Article) ListByCategory(categoryId int64, page, pageSize int) ([]model.Article, int64, error) {
	/*
	 * 1. 获取文章列表
	 * 2. 填充文章的分类信息
	 * 3. 查询文章总记录量
	 */
	fields := []string{"id", "category_id", "title", "views", "created_at", "updated_at"}

	articles, err := ctl.articleData.ListByCategory(fields, categoryId, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	if err := ctl.FillArticlesCategoryName(articles); err != nil {
		return nil, 0, err
	}

	totalSize, err := ctl.articleData.GetSumByCategory(categoryId)
	if err != nil {
		return nil, 0, err
	}

	return articles, totalSize, nil
}

func (ctl *Article) ListHome() ([]model.Article, error) {
	/*
	 * 1. 获取主页文章列表
	 * 2. 填充文章的分类信息
	 */
	fields := []string{"id", "category_id", "title", "preview_ctx", "views", "created_at", "updated_at"}
	const count = 5

	articles, err := ctl.articleData.ListByWeight(fields, count)
	if err != nil {
		return nil, err
	}

	if err := ctl.FillArticlesCategoryName(articles); err != nil {
		return nil, err
	}

	return articles, nil
}

// FillArticlesCategoryName 获取文章的分类名
func (ctl *Article) FillArticlesCategoryName(articles []model.Article) error {
	categories, err := ctl.categoryData.GetCategoriesMap()
	if err != nil {
		return err
	}

	for i := 0; i < len(articles); i++ {
		articles[i].CategoryName = categories[articles[i].CategoryId].Name
	}

	return nil
}
