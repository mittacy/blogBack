package service

import (
	"github.com/mittacy/blogBack/app/api"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/pkg/logger"
)

type Category struct {
	categoryData ICategoryData
	logger *logger.CustomLogger
}

// 编写实现api层中的各个service接口的构建方法

func NewCategory(categoryData ICategoryData, logger *logger.CustomLogger) api.ICategoryService {
	return &Category{
		categoryData: categoryData,
		logger: logger,
	}
}

type ICategoryData interface {
	Create(category *model.Category) error
	Delete(id int64) error
	UpdateNameById(category model.Category) error
	UpdateById(category model.Category, updateFields []string) error
	List(page, pageSize int) ([]model.Category, error)
	GetByName(name string) (*model.Category, error)
	GetSum() (int, error)
}

func (ctl *Category) Create(category model.Category) (int64, error) {
	if err := ctl.categoryData.Create(&category); err != nil {
		return 0, err
	}

	return category.Id, nil
}

func (ctl *Category) Delete(id int64) error {
	return ctl.categoryData.Delete(id)
}

func (ctl *Category) UpdateName(category model.Category) error {
	return ctl.categoryData.UpdateNameById(category)
}

func (ctl *Category) List(page, pageSize int) (categories []model.Category, totalSize int, err error) {
	// 查询列表
	if categories, err = ctl.categoryData.List(page, pageSize); err != nil {
		return
	}

	// 查询总记录数
	if totalSize, err = ctl.categoryData.GetSum(); err != nil {
		return
	}

	return categories, totalSize, nil
}

