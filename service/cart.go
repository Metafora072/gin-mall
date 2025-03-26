package service

import (
	"context"
	"gin-mail/dao"
	"gin-mail/model"
	"gin-mail/pkg/e"
	"gin-mail/serializer"
	"strconv"
)

type CartService struct {
	Id        uint `json:"cart_id" form:"cart_id"`
	BossId    uint `json:"boss_id" form:"boss_id"`
	ProductId uint `json:"product_id" form:"product_id"`
	Num       int  `json:"num" form:"num"`
}

func (service *CartService) Create(ctx context.Context, uid uint) serializer.Response {
	var cart *model.Cart
	code := e.Success

	// 先判断有没有相应的商品
	productDao := dao.NewProductDao(ctx)
	product, err := productDao.GetProductById(service.ProductId)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 判断有没有相应的卖家
	bossDao := dao.NewUserDao(ctx)
	_, err = bossDao.GetUserById(service.BossId)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 判断商品和卖家是否对应
	if product.BossId != service.BossId {
		code = e.ErrorProductAndBossNotMatch
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	cartDao := dao.NewCartDao(ctx)

	cart = &model.Cart{
		UserId:    uid,
		ProductId: service.ProductId,
		BossId:    service.BossId,
		Num:       uint(service.Num),
	}

	// 在数据库中插入 Cart 记录
	err = cartDao.CreateCart(cart)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func (service *CartService) List(ctx context.Context, uid uint) serializer.Response {
	code := e.Success

	cartDao := dao.NewCartDao(ctx)

	// 在数据库中获取用户的所有 Cart 记录
	cartList, err := cartDao.ListCartByUserId(uid)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildCarts(ctx, cartList),
	}
}

// Update 更新购物车，只考虑更新购买数量
func (service *CartService) Update(ctx context.Context, uid uint, cid string) serializer.Response {
	cartId, _ := strconv.Atoi(cid)
	code := e.Success

	cartDao := dao.NewCartDao(ctx)

	// 在数据库中更新 Cart 记录，只更新购买数量
	updated, err := cartDao.UpdateCartNumByIdAndUserId(uint(service.Num), uint(cartId), uid)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if !updated { // 没有匹配的 Cart 记录，修改失败
		code = e.ErrorCartNotFound
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func (service *CartService) Delete(ctx context.Context, uid uint, cid string) serializer.Response {
	cartId, _ := strconv.Atoi(cid)
	code := e.Success

	cartDao := dao.NewCartDao(ctx)

	// 在数据库中删除 Cart 记录
	deleted, err := cartDao.DeleteCartByIdAndUserId(uint(cartId), uid)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	if !deleted { // 没有匹配的 Cart 记录，删除失败
		code = e.ErrorCartNotFound
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
