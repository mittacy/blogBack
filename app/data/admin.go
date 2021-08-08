
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

type Admin struct {
	db 	   *gorm.DB
	logger *logger.CustomLogger
}

func NewAdmin(db *gorm.DB, logger *logger.CustomLogger) service.IAdminData {
	return &Admin{
		db:    	db,
		logger: logger,
	}
}

func (ctl *Admin) GetByName(name string) (*model.Admin, error) {
	var admin model.Admin
	if err := ctl.db.Where("name = ?", name).First(&admin).Error; err !=nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierr.ErrUserNoExist
		}
		return nil, errors.WithStack(err)
	}

	return &admin, nil
}

