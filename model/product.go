package model

import "gorm.io/gorm"

// Product 是描述商品的模型
// Name: 商品名称
// Category: 商品类型
// Title: 商品标题
// Info: 商品详情
// ImgPath: 商品展示图
// Price: 商品价格
// DiscountPrice: 商品打折后的价格
// OnSale: 商品是否在售
// Num: 商品数量
// BossId: 商品卖家 id
// BossName: 商品卖家名称
// BossAvatar: 商品卖家头像
type Product struct {
	gorm.Model
	Name          string
	Category      uint
	Title         string
	Info          string
	ImgPath       string
	Price         string
	DiscountPrice string
	OnSale        bool `gorm:"default:false"`
	Num           int
	BossId        uint
	BossName      string
	BossAvatar    string
}
