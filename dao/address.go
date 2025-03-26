package dao

import (
	"context"
	"gin-mail/model"
	"gorm.io/gorm"
)

type AddressDao struct {
	*gorm.DB
}

func NewAddressDao(ctx context.Context) *AddressDao {
	return &AddressDao{NewDBClient(ctx)}
}

func NewAddressDaoByDB(db *gorm.DB) *AddressDao {
	return &AddressDao{db}
}

func (dao *AddressDao) CreateAddress(address *model.Address) (err error) {
	err = dao.DB.Model(&model.Address{}).Create(&address).Error
	return
}

func (dao *AddressDao) GetAddressById(aid uint) (address *model.Address, err error) {
	err = dao.DB.Model(&model.Address{}).Where("id = ?", aid).First(&address).Error
	return
}

// ListAddressByUserId 获取用户的所有 address 记录
func (dao *AddressDao) ListAddressByUserId(uid uint) (addresses []*model.Address, err error) {
	err = dao.DB.Model(&model.Address{}).Where("user_id = ?", uid).Find(&addresses).Error
	return
}

// UpdateAddressByIdAndUserId 根据 address id 和用户 id 修改相应的 address 记录
// updated: 是否成功修改
func (dao *AddressDao) UpdateAddressByIdAndUserId(aid uint, uid uint, address *model.Address) (updated bool, err error) {
	result := dao.DB.Model(&model.Address{}).Where("id = ? AND user_id = ?", aid, uid).Updates(&address)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

// DeleteAddressByIdAndUserId 根据 address id 和用户 id 删除相应的 address 记录
// deleted: 是否成功删除
func (dao *AddressDao) DeleteAddressByIdAndUserId(aid, uid uint) (deleted bool, err error) {
	result := dao.DB.Model(&model.Address{}).Where("id = ? AND user_id = ?", aid, uid).Delete(&model.Address{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
