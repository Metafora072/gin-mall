package service

import (
	"context"
	"gin-mail/dao"
	"gin-mail/pkg/utils"
	"gin-mail/serializer"
	"strconv"
)

type ListProductImg struct {
}

func (service *ListProductImg) List(ctx context.Context, pid string) serializer.Response {
	productImgDao := dao.NewProductImgDao(ctx)
	productId, _ := strconv.Atoi(pid)

	// 获取指定商品 id 对应的若干个图片
	productImgs, err := productImgDao.ListProductImg(uint(productId))
	if err != nil {
		utils.LogrusObj.Infoln("ListProductImg func List err: ", err)
	}

	return serializer.BuildListResponse(serializer.BuildProductImgs(productImgs), uint(len(productImgs)))
}
