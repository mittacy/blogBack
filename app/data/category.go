package data

import (
	"github.com/mittacy/blogBack/apierr"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/app/service"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// 实现service层中的data接口

// 分类缓存较少，直接使用内存缓存
type categoryCache struct {
	m       map[int64]model.Category // 分类map
	s       []model.Category         // 所有分类切片
	isValid bool                     // 是否有效
}

var categoryData categoryCache // 缓存所有分类

func cacheCategories(categories []model.Category) {
	categoryData.s = categories
	categoryData.m = make(map[int64]model.Category, len(categoryData.s))
	for _, v := range categoryData.s {
		categoryData.m[v.Id] = v
	}
	categoryData.isValid = true
}

func init() {
	categoryData = categoryCache{
		m:       make(map[int64]model.Category, 0),
		s:       make([]model.Category, 0),
		isValid: false,
	}
}

type Category struct {
	db     *gorm.DB
	logger *logger.CustomLogger
}

func NewCategory(db *gorm.DB, logger *logger.CustomLogger) service.ICategoryData {
	return &Category{
		db:     db,
		logger: logger,
	}
}

func NewArticleCategory(db *gorm.DB, logger *logger.CustomLogger) service.IArticleCategoryData {
	return &Category{
		db:     db,
		logger: logger,
	}
}

func (ctl *Category) Create(category *model.Category) error {
	// 查询name是否存在
	count := 0
	err := ctl.db.Model(category).Select("1").Where("name = ?", category.Name).First(&count).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if count > 0 {
		return apierr.ErrCategoryNameExist
	}

	// 创建分类
	if err := ctl.db.Create(category).Error; err != nil {
		return err
	}

	ctl.ExpireCategoryData()

	return nil
}

func (ctl *Category) Delete(id int64) error {
	// 从数据库删除
	category := model.Category{Id: id}
	res := ctl.db.Delete(&category)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return apierr.ErrCategoryNoExist
	}

	ctl.ExpireCategoryData()
	return nil
}

func (ctl *Category) UpdateNameById(category model.Category) error {
	// 查询name是否存在
	existCategory := model.Category{}
	err := ctl.db.Model(category).Select("id").Where("name = ?", category.Name).First(&existCategory).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) && existCategory.Id != category.Id {
		return apierr.ErrCategoryNameExist
	}

	// 更新
	return ctl.UpdateById(category, []string{"name"})

}

func (ctl *Category) UpdateById(category model.Category, updateFields []string) error {
	if err := ctl.db.Select(updateFields).Updates(&category).Error; err != nil {
		return errors.WithStack(err)
	}

	ctl.ExpireCategoryData()
	return nil
}

func (ctl *Category) List(page, pageSize int) (categories []model.Category, err error) {
	/*
	 * 1. 不存在 -> 数据库查询, 存入缓存
	 * 2. 分页并返回
	 */
	if !categoryData.isValid {
		if err = ctl.db.Find(&categories).Error; err != nil {
			return nil, errors.WithStack(err)
		}

		// 缓存
		cacheCategories(categories)
	}

	// 不分页，返回全部
	if pageSize == 0 {
		return categoryData.s, nil
	}

	// 分页返回
	return ctl.dataPage(categoryData.s, page, pageSize), nil
}

func (ctl *Category) GetCategoriesMap() (map[int64]model.Category, error) {
	if !categoryData.isValid {
		if _, err := ctl.List(0, 0); err != nil {
			return nil, err
		}
	}
	return categoryData.m, nil
}

func (ctl *Category) GetSum() (int, error) {
	if !categoryData.isValid {
		if _, err := ctl.List(0, 0); err != nil {
			return 0, err
		}
	}

	return len(categoryData.s), nil
}

func (ctl *Category) GetByName(name string) (*model.Category, error) {
	category := model.Category{}

	if err := ctl.db.Where("name = ?", name).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

func (ctl *Category) ExpireCategoryData() {
	categoryData.isValid = false
}

func (ctl *Category) dataPage(categories []model.Category, page, pageSize int) []model.Category {
	// 不分页
	if pageSize == 0 || len(categories) == 0 {
		return categories
	}

	// 分页
	if pageSize <= 0 {
		pageSize = 2
	}
	if page <= 0 {
		page = 1
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(categories) {
		return []model.Category{}
	}
	if end > len(categories) {
		end = len(categories)
	}

	return categories[start:end]
}
