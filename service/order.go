package service

import (
	"context"
	"fmt"
	"gin-mail/dao"
	"gin-mail/model"
	"gin-mail/pkg/e"
	"gin-mail/serializer"
	"math/rand"
	"strconv"
	"time"
)

type OrderService struct {
	ProductId uint    `json:"product_id" form:"product_id"`
	Num       int     `json:"num" form:"num"`
	AddressId uint    `json:"address_id" form:"address_id"`
	Money     float64 `json:"money" form:"money"`
	BossId    uint    `json:"boss_id" form:"boss_id"`
	UserId    uint    `json:"user_id" form:"user_id"`
	OrderNum  int     `json:"order_num" form:"order_num"`
	Type      int     `json:"type" form:"type"`
	model.BasePage
}

func (service *OrderService) Create(ctx context.Context, uid uint) serializer.Response {
	var order *model.Order
	code := e.Success

	// 先判断有没有相应的商品
	productDao := dao.NewProductDao(ctx)
	product, err := productDao.GetProductById(service.ProductId)
	if err != nil {
		code = e.ErrorOrderProductNotFound
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
		code = e.ErrorOrderBossNotFound
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

	// 检验订单地址是否存在
	addressDao := dao.NewAddressDao(ctx)
	address, err := addressDao.GetAddressById(service.AddressId)
	if err != nil {
		code = e.ErrorOrderAddressNotFound
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	orderDao := dao.NewOrderDao(ctx)

	order = &model.Order{
		UserId:    uid,
		ProductId: service.ProductId,
		BossId:    service.BossId,
		Num:       service.Num,
		Money:     service.Money,
		Type:      model.OrderTypeNotPaid, // 默认订单未支付
		AddressId: address.ID,
	}

	// 生成唯一订单号赋值给 OrderNum
	// 订单号： 自动生成的随机 number + product_id + user_id
	number := fmt.Sprintf("%09v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	productNum := strconv.Itoa(int(service.ProductId))
	userNum := strconv.Itoa(int(service.UserId))
	number = number + productNum + userNum
	order.OrderNum, _ = strconv.ParseUint(number, 10, 64)

	// 在数据库中插入 Order 记录
	err = orderDao.CreateOrder(order)
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

// Show 获取相应 id 订单的详细信息
func (service *OrderService) Show(ctx context.Context, uid uint, oid string) serializer.Response {
	orderId, _ := strconv.Atoi(oid)
	code := e.Success

	orderDao := dao.NewOrderDao(ctx)

	// 在数据库中获取 order 记录
	order, err := orderDao.GetOrderByIdAndUserId(uint(orderId), uid)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 获取订单的地址信息
	addressDao := dao.NewAddressDao(ctx)
	address, err := addressDao.GetAddressById(order.AddressId)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 获取订单的商品信息
	productDao := dao.NewProductDao(ctx)
	product, err := productDao.GetProductById(order.ProductId)
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
		Data:   serializer.BuildOrder(order, product, address),
	}
}

func (service *OrderService) List(ctx context.Context, uid uint) serializer.Response {
	code := e.Success

	if service.PageSize == 0 {
		service.PageSize = 15
	}

	// 指定查询条件 condition
	condition := make(map[string]interface{})
	if service.Type != 0 { // 查询特定 Type
		condition["type"] = service.Type
	}
	condition["user_id"] = uid

	orderDao := dao.NewOrderDao(ctx)
	// 在数据库中获取用户的所有 Order 记录
	orderList, total, err := orderDao.ListOrderByCondition(condition, service.BasePage)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.BuildListResponse(serializer.BuildOrders(ctx, orderList), uint(total))
}

func (service *OrderService) Delete(ctx context.Context, uid uint, oid string) serializer.Response {
	orderId, _ := strconv.Atoi(oid)
	code := e.Success

	orderDao := dao.NewOrderDao(ctx)

	// 在数据库中删除 Order 记录
	deleted, err := orderDao.DeleteOrderByIdAndUserId(uint(orderId), uid)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	if !deleted { // 没有匹配的 Order 记录，删除失败
		code = e.ErrorOrderNotFound
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
