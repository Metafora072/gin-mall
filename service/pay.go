package service

import (
	"context"
	"errors"
	"fmt"
	"gin-mail/dao"
	"gin-mail/pkg/e"
	"gin-mail/pkg/utils"
	"gin-mail/serializer"
	"gorm.io/gorm"
	"strconv"
)

type OrderPay struct {
	OrderId   uint    `json:"order_id" form:"order_id"`
	Money     float64 `json:"money" form:"money"`
	OrderNo   string  `json:"order_no" form:"order_no"`
	ProductId uint    `json:"product_id" form:"product_id"`
	PayTime   string  `json:"pay_time" form:"pay_time"`
	Sign      string  `json:"sign" form:"sign"`
	BossId    uint    `json:"boss_id" form:"boss_id"`
	BossName  string  `json:"boss_name" form:"boss_name"`
	Num       int     `json:"num" form:"num"`
	Key       string  `json:"key" form:"key"`
}

func (service *OrderPay) Pay(ctx context.Context, uid uint) serializer.Response {
	utils.Encrypt.SetKey(service.Key)
	code := e.Success
	err := dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		orderDao := dao.NewOrderDaoByDB(tx)

		order, err := orderDao.GetOrderById(service.OrderId)
		if err != nil {
			utils.LogrusObj.Infoln("OrderPay func Pay: ", err)
			code = e.Error
			return err
		}

		// 订单的价格 money
		money := order.Money * float64(order.Num)

		// 订单买家 user
		userDao := dao.NewUserDaoByDB(tx)
		user, err := userDao.GetUserById(uid)
		if err != nil {
			utils.LogrusObj.Infoln("OrderPay func Pay: ", err)
			code = e.Error
			return err
		}

		// 对 user 的钱进行解密，减去购买金额，再重新加密回去
		userMoneyStr := utils.Encrypt.AesDecoding(user.Money)
		userMoneyFloat, _ := strconv.ParseFloat(userMoneyStr, 64)

		if userMoneyFloat < money { // 用户余额不足
			code = e.Error
			utils.LogrusObj.Infoln("OrderPay func Pay: ", errors.New("用户余额不足"))
			return errors.New("用户余额不足")
		}

		// 计算用户扣除订单价格的余额
		userFinalMoney := fmt.Sprintf("%f", userMoneyFloat-money)
		user.Money = utils.Encrypt.AesEncoding(userFinalMoney)

		userDao = dao.NewUserDaoByDB(userDao.DB)
		err = userDao.UpdateUserById(uid, user)
		if err != nil { // 更新用户余额失败
			utils.LogrusObj.Infoln("OrderPay func Pay: ", err)
			code = e.Error
			return err
		}

		// 订单卖家 boss
		boss, err := userDao.GetUserById(order.BossId)
		if err != nil {
			utils.LogrusObj.Infoln("OrderPay func Pay: ", err)
			code = e.Error
			return err
		}

		// 对 boss 的钱进行解密，加上购买金额，再重新加密回去
		bossMoneyStr := utils.Encrypt.AesDecoding(boss.Money)
		bossMoneyFloat, _ := strconv.ParseFloat(bossMoneyStr, 64)

		// 计算卖家加上订单价格的余额
		bossFinalMoney := fmt.Sprintf("%f", bossMoneyFloat+money)
		boss.Money = utils.Encrypt.AesEncoding(bossFinalMoney)

		err = userDao.UpdateUserById(boss.ID, boss)
		if err != nil {
			utils.LogrusObj.Infoln("OrderPay func Pay: ", err)
			code = e.Error
			return err
		}

		// 更新商品数量
		productDao := dao.NewProductDaoByDB(tx)
		product, err := productDao.GetProductById(order.ProductId)
		if err != nil {
			utils.LogrusObj.Infoln("OrderPay func Pay: ", err)
			code = e.Error
			return err
		}

		product.Num -= order.Num
		err = productDao.UpdateProductById(order.ProductId, product)
		if err != nil {
			utils.LogrusObj.Infoln("OrderPay func Pay: ", err)
			code = e.Error
			return err
		}

		// 删除这个 order 订单
		deleted, err := orderDao.DeleteOrderByIdAndUserId(order.ID, uid)
		if err != nil {
			utils.LogrusObj.Infoln("OrderPay func Pay: ", err)
			code = e.Error
			return err
		}
		if !deleted {
			utils.LogrusObj.Infoln("OrderPay func Pay: ", errors.New("订单删除失败"))
			code = e.Error
			return errors.New("订单删除失败")
		}

		return nil
	})
	if err != nil {
		code = e.Error
		utils.LogrusObj.Infoln("OrderPay func Pay: ", err)
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
