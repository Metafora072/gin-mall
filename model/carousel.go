package model

import "gorm.io/gorm"

// Carousel 描述轮播图
type Carousel struct {
	gorm.Model
	ImgPath   string
	ProductId uint `gorm:"not null"`
}
