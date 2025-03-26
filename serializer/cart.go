package serializer

import (
	"context"
	"gin-mail/conf"
	"gin-mail/dao"
	"gin-mail/model"
	"gin-mail/pkg/utils"
)

type CartVO struct {
	Id            uint   `json:"id"`
	UserId        uint   `json:"user_id"`
	ProductId     uint   `json:"product_id"`
	CreatedAt     int64  `json:"created_at"`
	Name          string `json:"name"`
	Num           uint   `json:"num"`
	MaxNum        uint   `json:"max_num"`
	ImgPath       string `json:"img_path"`
	Check         bool   `json:"check"`
	Price         string `json:"price"`
	DiscountPrice string `json:"discount_price"`
	BossId        uint   `json:"boss_id"`
	BossName      string `json:"boss_name"`
}

func BuildCart(cart *model.Cart, product *model.Product) CartVO {
	return CartVO{
		Id:            cart.ID,
		UserId:        cart.UserId,
		ProductId:     cart.ProductId,
		CreatedAt:     cart.CreatedAt.Unix(),
		Name:          product.Name,
		Num:           cart.Num,
		MaxNum:        cart.MaxNum,
		ImgPath:       conf.Host + conf.HttpPort + conf.ProductPath + product.ImgPath,
		Check:         cart.Check,
		Price:         product.Price,
		DiscountPrice: product.DiscountPrice,
		BossId:        product.BossId,
		BossName:      product.BossName,
	}
}

func BuildCarts(ctx context.Context, items []*model.Cart) (carts []CartVO) {
	productDao := dao.NewProductDao(ctx)
	for _, item := range items {
		product, err := productDao.GetProductById(item.ProductId)
		if err != nil {
			utils.LogrusObj.Infoln("CartVO func BuildCarts: ", err)
			continue
		}
		carts = append(carts, BuildCart(item, product))
	}
	return
}
