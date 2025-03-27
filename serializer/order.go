package serializer

import (
	"context"
	"gin-mail/conf"
	"gin-mail/dao"
	"gin-mail/model"
	"gin-mail/pkg/utils"
)

type OrderVO struct {
	Id           uint    `json:"id"`
	OrderNum     uint64  `json:"order_num"`
	CreatedAt    int64   `json:"created_at"`
	UpdatedAt    int64   `json:"updated_at"`
	UserId       uint    `json:"user_id"`
	ProductId    uint    `json:"product_id"`
	BossId       uint    `json:"boss_id"`
	Num          int     `json:"num"`
	Address      string  `json:"address"`
	AddressName  string  `json:"address_name"`
	AddressPhone string  `json:"address_phone"`
	Type         uint    `json:"type"`
	ProductName  string  `json:"product_name"`
	ImgPath      string  `json:"img_path"`
	Money        float64 `json:"money"`
}

func BuildOrder(order *model.Order, product *model.Product, address *model.Address) OrderVO {
	return OrderVO{
		Id:           order.ID,
		OrderNum:     order.OrderNum,
		CreatedAt:    order.CreatedAt.Unix(),
		UpdatedAt:    order.UpdatedAt.Unix(),
		UserId:       order.UserId,
		ProductId:    order.ProductId,
		BossId:       order.BossId,
		Num:          order.Num,
		Address:      address.Address,
		AddressName:  address.Name,
		AddressPhone: address.Phone,
		Type:         order.Type,
		ProductName:  product.Name,
		ImgPath:      conf.Host + conf.HttpPort + conf.ProductPath + product.ImgPath,
		Money:        order.Money,
	}
}

func BuildOrders(ctx context.Context, items []*model.Order) (orders []OrderVO) {
	productDao := dao.NewProductDao(ctx)
	addressDao := dao.NewAddressDao(ctx)

	for _, item := range items {
		product, err := productDao.GetProductById(item.ProductId)
		if err != nil {
			utils.LogrusObj.Infoln("OrderVO func BuildFavorites: ", err)
			continue
		}
		address, err := addressDao.GetAddressById(item.AddressId)
		if err != nil {
			utils.LogrusObj.Infoln("OrderVO func BuildFavorites: ", err)
			continue
		}
		orders = append(orders, BuildOrder(item, product, address))
	}
	return
}
