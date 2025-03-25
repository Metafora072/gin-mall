package service

import (
	"context"
	"gin-mail/dao"
	"gin-mail/model"
	"gin-mail/pkg/e"
	"gin-mail/pkg/utils"
	"gin-mail/serializer"
	"strconv"
)

// FavoriteService 是收藏夹绑定的结构体
// FavoriteId: 收藏夹标识 id
// ProductId: 收藏夹记录的商品 id
// BossId: 收藏夹记录商品的卖家 id
// model.BasePage: 分页模式
type FavoriteService struct {
	FavoriteId uint `json:"favorite_id" form:"favorite_id"`
	ProductId  uint `json:"product_id" form:"product_id"`
	BossId     uint `json:"boss_id" form:"boss_id"`
	model.BasePage
}

// List 获取收藏夹信息
func (service *FavoriteService) List(ctx context.Context, uid uint) serializer.Response {
	favoriteDao := dao.NewFavoriteDao(ctx)
	code := e.Success
	favorites, err := favoriteDao.ListFavorites(uid)
	if err != nil {
		utils.LogrusObj.Infoln("Favorite Service func List: ", err)
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.BuildListResponse(serializer.BuildFavorites(ctx, favorites), uint(len(favorites)))
}

// Create 创建收藏夹
func (service *FavoriteService) Create(ctx context.Context, uid uint) serializer.Response {
	code := e.Success
	favoriteDao := dao.NewFavoriteDao(ctx)
	// 检查创建的收藏夹 id 是否已经存在
	exist, _ := favoriteDao.FavoriteExistOrNot(service.ProductId, uid)
	if exist {
		code = e.ErrorFavoriteExist // 收藏夹记录已经存在
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  e.GetMsg(code),
		}
	}

	// 获取收藏夹对应的用户 user 信息
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(uid)
	if err != nil {
		code = e.Error
		utils.LogrusObj.Infoln("Favorite Service func Create in userDao.GetUserById: ", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 获取收藏夹对应商品的卖家 boss 信息
	bossDao := dao.NewUserDao(ctx)
	boss, err := bossDao.GetUserById(service.BossId)
	if err != nil {
		code = e.Error
		utils.LogrusObj.Infoln("Favorite Service func Create in bossDao.GetUserById: ", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 获取收藏夹对应的商品 product 信息
	productDao := dao.NewProductDao(ctx)
	product, err := productDao.GetProductById(service.ProductId)
	if err != nil {
		code = e.Error
		utils.LogrusObj.Infoln("Favorite Service func Create in productDao.GetProductById: ", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 检查商品是否和商家匹配
	// TODO: 为优化性能，尽量在 productDao.GetProductById 将 boss.Id 纳入查询条件中 (添加外键)
	if product.BossId != boss.ID {
		code = e.Error
		utils.LogrusObj.Infoln("Favorite Service Create: product id != boss id")
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "商品和商家不匹配",
		}
	}

	// 创建收藏夹结构体 favorite
	favorite := &model.Favorite{
		User:      *user,
		UserId:    uid,
		Product:   *product,
		ProductId: service.ProductId,
		Boss:      *boss,
		BossId:    service.BossId,
	}

	// 将 favorite 写入数据库
	err = favoriteDao.CreateFavorite(favorite)
	if err != nil {
		code = e.Error
		utils.LogrusObj.Infoln("Favorite Service func Create in favoriteDao.CreateFavorite: ", err)
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

// Delete 删除收藏夹
func (service *FavoriteService) Delete(ctx context.Context, uid uint, fid string) serializer.Response {
	favoriteDao := dao.NewFavoriteDao(ctx)
	favoriteId, _ := strconv.Atoi(fid)
	utils.LogrusObj.Infoln("Favorite Service Delete: ", favoriteId)
	code := e.Success
	// 在数据库里是软删除，即添加 deleted_at 字段
	deleted, err := favoriteDao.DeleteFavorite(uid, uint(favoriteId))
	if err != nil {
		code = e.Error
		utils.LogrusObj.Infoln("Favorite Service func Delete in favoriteDao.DeleteFavorite: ", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if !deleted {
		code = e.ErrorFavoriteDeleteNotFound
		utils.LogrusObj.Infoln("Favorite Service func Delete in favoriteDao.DeleteFavorite: ", e.ErrorFavoriteDeleteNotFound)
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
