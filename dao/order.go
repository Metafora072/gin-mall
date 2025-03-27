package dao

import (
	"context"
	"gin-mail/model"
	"gorm.io/gorm"
)

type OrderDao struct {
	*gorm.DB
}

func NewOrderDao(ctx context.Context) *OrderDao {
	return &OrderDao{NewDBClient(ctx)}
}

func NewOrderDaoByDB(db *gorm.DB) *OrderDao {
	return &OrderDao{db}
}

func (dao *OrderDao) CreateOrder(order *model.Order) (err error) {
	err = dao.DB.Model(&model.Order{}).Create(&order).Error
	return
}

func (dao *OrderDao) GetOrderById(oid uint) (order *model.Order, err error) {
	err = dao.DB.Model(&model.Order{}).Where("id = ?", oid).First(&order).Error
	return
}

func (dao *OrderDao) GetOrderByIdAndUserId(oid, uid uint) (order *model.Order, err error) {
	err = dao.DB.Model(&model.Order{}).Where("id = ? AND user_id = ?", oid, uid).First(&order).Error
	return
}

// ListOrderByUserId 获取用户的所有 Order 记录
func (dao *OrderDao) ListOrderByUserId(uid uint) (orders []*model.Order, err error) {
	err = dao.DB.Model(&model.Order{}).Where("user_id = ?", uid).Find(&orders).Error
	return
}

// UpdateOrderByIdAndUserId 根据 Order id 和用户 id 修改相应的 Order 记录
// updated: 是否成功修改
func (dao *OrderDao) UpdateOrderByIdAndUserId(oid uint, uid uint, order *model.Order) (updated bool, err error) {
	result := dao.DB.Model(&model.Order{}).Where("id = ? AND user_id = ?", oid, uid).Updates(&order)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

// DeleteOrderByIdAndUserId 根据 Order id 和用户 id 删除相应的 Order 记录
// deleted: 是否成功删除
func (dao *OrderDao) DeleteOrderByIdAndUserId(oid, uid uint) (deleted bool, err error) {
	result := dao.DB.Model(&model.Order{}).Where("id = ? AND user_id = ?", oid, uid).Delete(&model.Order{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

// ListOrderByCondition 根据查询条件 condition 获取符合的所有 order 记录，并以指定分页模式分页
func (dao *OrderDao) ListOrderByCondition(condition map[string]interface{}, page model.BasePage) (orders []*model.Order, total int64, err error) {
	err = dao.DB.Model(&model.Order{}).Where(condition).Count(&total).Error
	if err != nil {
		return
	}
	err = dao.DB.Model(&model.Order{}).Where(condition).Offset((page.PageNum - 1) * page.PageSize).Limit(page.PageSize).Find(&orders).Error
	return
}
