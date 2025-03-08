package model

import (
	"gorm.io/gorm"
)

// Cart 是描述购物车的模型
// UserId: 购物车对应的用户 id
// ProductId: 商品 id
// BossId: 商品卖家 id
// Num:	购买数量
// MaxNum: 最大购买数量
// Check: 是否已支付
type Cart struct {
	gorm.Model
	UserId    uint `gorm:"not null"`
	ProductId uint `gorm:"not null"`
	BossId    uint `gorm:"not null"`
	Num       uint `gorm:"not null"`
	MaxNum    uint `gorm:"not null"`
	Check     bool
}
