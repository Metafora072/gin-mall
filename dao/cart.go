package dao

import (
	"context"
	"gin-mail/model"
	"gorm.io/gorm"
)

type CartDao struct {
	*gorm.DB
}

func NewCartDao(ctx context.Context) *CartDao {
	return &CartDao{NewDBClient(ctx)}
}

func NewCartDaoByDB(db *gorm.DB) *CartDao {
	return &CartDao{db}
}

func (dao *CartDao) CreateCart(cart *model.Cart) (err error) {
	err = dao.DB.Model(&model.Cart{}).Create(&cart).Error
	return
}

func (dao *CartDao) GetCartById(cid uint) (cart *model.Cart, err error) {
	err = dao.DB.Model(&model.Cart{}).Where("id = ?", cid).First(&cart).Error
	return
}

// ListCartByUserId 获取用户的所有 Cart 记录
func (dao *CartDao) ListCartByUserId(uid uint) (carts []*model.Cart, err error) {
	err = dao.DB.Model(&model.Cart{}).Where("user_id = ?", uid).Find(&carts).Error
	return
}

// UpdateCartByIdAndUserId 根据 Cart id 和用户 id 修改相应的 Cart 记录
// updated: 是否成功修改
func (dao *CartDao) UpdateCartByIdAndUserId(cid uint, uid uint, cart *model.Cart) (updated bool, err error) {
	result := dao.DB.Model(&model.Cart{}).Where("id = ? AND user_id = ?", cid, uid).Updates(&cart)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

// DeleteCartByIdAndUserId 根据 Cart id 和用户 id 删除相应的 Cart 记录
// deleted: 是否成功删除
func (dao *CartDao) DeleteCartByIdAndUserId(cid, uid uint) (deleted bool, err error) {
	result := dao.DB.Model(&model.Cart{}).Where("id = ? AND user_id = ?", cid, uid).Delete(&model.Cart{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (dao *CartDao) UpdateCartNumByIdAndUserId(num, cid, uid uint) (updated bool, err error) {
	result := dao.DB.Model(&model.Cart{}).Where("id = ? AND user_id = ?", cid, uid).Update("num", num)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
