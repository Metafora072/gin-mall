package serializer

import (
	"context"
	"gin-mail/conf"
	"gin-mail/model"
)

type ProductVO struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	CategoryId    uint   `json:"category_id"`
	Title         string `json:"title"`
	Info          string `json:"info"`
	ImgPath       string `json:"img_path"`
	Price         string `json:"price"`
	DiscountPrice string `json:"discount_price"`
	View          uint64 `json:"view"`
	CreatedAt     int64  `json:"created_at"`
	Num           int    `json:"num"`
	OnSale        bool   `json:"on_sale"`
	BossId        uint   `json:"boss_id"`
	BossName      string `json:"boss_name"`
	BossAvatar    string `json:"boss_avatar"`
}

func BuildProduct(ctx context.Context, item *model.Product) ProductVO {
	return ProductVO{
		Id:            item.ID,
		Name:          item.Name,
		CategoryId:    item.CategoryId,
		Title:         item.Title,
		Info:          item.Info,
		ImgPath:       conf.Host + conf.HttpPort + conf.ProductPath + item.ImgPath,
		Price:         item.Price,
		DiscountPrice: item.DiscountPrice,
		View:          item.View(ctx),
		CreatedAt:     item.CreatedAt.Unix(),
		Num:           item.Num,
		OnSale:        item.OnSale,
		BossId:        item.BossId,
		BossName:      item.BossName,
		BossAvatar:    conf.Host + conf.HttpPort + conf.AvatarPath + item.BossAvatar,
	}

}

func BuildProducts(ctx context.Context, items []*model.Product) (products []ProductVO) {
	for _, item := range items {
		product := BuildProduct(ctx, item)
		products = append(products, product)
	}
	return
}
