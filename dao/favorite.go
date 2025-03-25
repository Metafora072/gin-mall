package dao

import (
	"context"
	"errors"
	"gin-mail/model"
	"gorm.io/gorm"
)

type FavoriteDao struct {
	*gorm.DB
}

func NewFavoriteDao(ctx context.Context) *FavoriteDao {
	return &FavoriteDao{NewDBClient(ctx)}
}

func NewFavoriteDaoByDB(db *gorm.DB) *FavoriteDao {
	return &FavoriteDao{db}
}

// ListFavorites 根据用户 uid 获取其所有收藏夹信息
func (dao *FavoriteDao) ListFavorites(uid uint) (favorites []*model.Favorite, err error) {
	err = dao.DB.Model(&model.Favorite{}).Where("user_id = ?", uid).Find(&favorites).Error
	return
}

// FavoriteExistOrNot 根据商品 pid 和用户 uid 检查数据库中是否存在相应记录
func (dao *FavoriteDao) FavoriteExistOrNot(pid, uid uint) (exist bool, err error) {
	var count int64
	err = dao.DB.Model(&model.Favorite{}).Where("user_id = ? AND product_id = ?", uid, pid).Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 没有相应记录
			return false, nil
		}
		return false, err // 数据库查询错误
	}
	return count > 0, nil
}

// CreateFavorite 在数据库中创建一条 favorite 收藏夹记录
func (dao *FavoriteDao) CreateFavorite(favorite *model.Favorite) (err error) {
	err = dao.DB.Model(&model.Favorite{}).Create(&favorite).Error
	return
}

// DeleteFavorite 在数据库中根据用户 uid 和收藏夹 id 删除相应的收藏夹记录
// deleted: 是否成功删除
func (dao *FavoriteDao) DeleteFavorite(uid, fid uint) (deleted bool, err error) {
	result := dao.DB.Model(&model.Favorite{}).Where("id = ? AND user_id = ?", fid, uid).Delete(&model.Favorite{})
	if result.Error != nil {
		return false, result.Error
	}
	// result.RowsAffected 代表影响的行数，也就是删除的行数，如果大于 0 代表成功删除，否则就是数据库中没有相应记录
	return result.RowsAffected > 0, nil
}
