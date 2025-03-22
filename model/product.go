package model

import (
	"context"
	"gin-mail/cache"
	"gorm.io/gorm"
	"strconv"
)

// Product 是描述商品的模型
// Name: 商品名称
// CategoryId: 商品类型
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
	CategoryId    uint
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

// View 获取商品点击数
func (product *Product) View(ctx context.Context) uint64 {
	// 从 Redis 获取指定商品的点击数，键名由 cache.ProductViewKey(product.ID) 生成
	countStr, _ := cache.RedisClient.Get(ctx, cache.ProductViewKey(product.ID)).Result()
	// 将获取到的字符串类型的点击数转换为 uint64 类型
	count, _ := strconv.ParseUint(countStr, 10, 64)
	return count
}

// AddView 增加商品点击数
func (product *Product) AddView(ctx context.Context) {
	// 在 Redis 中将指定商品的点击数自增 1
	cache.RedisClient.Incr(ctx, cache.ProductViewKey(product.ID))
	// 增加商品排名
	// 在 Redis 的有序集合 (ZSet) 中，为指定商品的排名值自增 1, ZSet 是 Redis 中的有序集合，自动按分数排序，方便实现排行榜功能。
	// cache.RedisClient.ZIncrBy()：
	// Redis 的 ZINCRBY 命令，增加指定商品在有序集合 (ZSet) 中的分数（用于排名管理）。
	// cache.RankKey 是存储商品排名的键名。
	// 第三个参数 strconv.Itoa(int(product.ID)) 将商品 ID 转换为字符串形式，作为 ZSet 的成员。
	// ZSet Key: "rank"
	// -----------------
	// Member | Score
	// -------|-------
	// 1001   | 15     // 商品ID 1001 点击数15次
	// 1002   | 25     // 商品ID 1002 点击数25次
	// 1003   | 10     // 商品ID 1003 点击数10次
	cache.RedisClient.ZIncrBy(ctx, cache.RankKey, 1, strconv.Itoa(int(product.ID)))
}
