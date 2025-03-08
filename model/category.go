package model

import "gorm.io/gorm"

// Category 是描述商品分类的模型
type Category struct {
	gorm.Model
	CategoryName string
}
