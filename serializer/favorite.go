package serializer

import (
	"context"
	"gin-mail/conf"
	"gin-mail/dao"
	"gin-mail/model"
	"gin-mail/pkg/utils"
)

type FavoriteVO struct {
	UserId        uint   `json:"user_id"`
	ProductId     uint   `json:"product_id"`
	CreatedAt     int64  `json:"created_at"`
	Name          string `json:"name"`
	CategoryId    uint   `json:"category_id"`
	Title         string `json:"title"`
	Info          string `json:"info"`
	ImgPath       string `json:"img_path"`
	Price         string `json:"price"`
	DiscountPrice string `json:"discount_price"`
	BossId        uint   `json:"boss_id"`
	Num           int    `json:"num"`
	OnSale        bool   `json:"on_sale"`
}

func BuildFavorite(favorite *model.Favorite, product *model.Product) FavoriteVO {
	return FavoriteVO{
		UserId:        favorite.UserId,
		ProductId:     favorite.ProductId,
		CreatedAt:     favorite.CreatedAt.Unix(),
		Name:          product.Name,
		CategoryId:    product.CategoryId,
		Title:         product.Title,
		Info:          product.Info,
		ImgPath:       conf.Host + conf.HttpPort + conf.ProductPath + product.ImgPath,
		Price:         product.Price,
		DiscountPrice: product.DiscountPrice,
		BossId:        product.BossId,
		Num:           product.Num,
		OnSale:        product.OnSale,
	}
}

func BuildFavorites(ctx context.Context, items []*model.Favorite) (favorites []FavoriteVO) {
	productDao := dao.NewProductDao(ctx)
	for _, item := range items {
		product, err := productDao.GetProductById(item.ProductId)
		if err != nil {
			utils.LogrusObj.Infoln("FavoriteVO func BuildFavorites: ", err)
			continue
		}
		favorites = append(favorites, BuildFavorite(item, product))
	}
	return
}
