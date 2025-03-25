package dao

import (
	"context"
	"gin-mail/model"
	"gorm.io/gorm"
)

type CategoryDao struct {
	*gorm.DB
}

func NewCategoryDao(ctx context.Context) *CategoryDao {
	return &CategoryDao{NewDBClient(ctx)}
}

func NewCategoryDaoByDB(db *gorm.DB) *CategoryDao {
	return &CategoryDao{db}
}

// ListCategory 列举商品分类
func (dao *CategoryDao) ListCategory() (Category []*model.Category, err error) {
	err = dao.DB.Model(&model.Category{}).Find(&Category).Error
	return
}
