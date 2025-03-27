package model

import "gorm.io/gorm"

// Order 是描述订单的模型
// UserId: 订单用户 id
// ProductId: 订单商品 id
// BossId: 订单商品的卖家 id
// AddressId: 订单地址 id
// Num: 数量
// OrderNum: 订单号
// Type: 订单类型。1表示未支付，2表示已支付
// Money: 订单价格
type Order struct {
	gorm.Model
	UserId    uint `gorm:"not null"`
	ProductId uint `gorm:"not null"`
	BossId    uint `gorm:"not null"`
	AddressId uint `gorm:"not null"`
	Num       int
	OrderNum  uint64
	Type      uint
	Money     float64
}

// 描述 Order.Type 字段(订单是否已支付)
const (
	OrderTypeNotPaid = 1
	OrderTypePaid    = 2
)
