package model

import (
	"gorm.io/gorm"
)

// ProductImg 是描述商品图片的模型
type ProductImg struct {
	gorm.Model
	ProductId uint `gorm:"not null"`
	ImgPath   string
}
