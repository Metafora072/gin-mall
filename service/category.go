package service

import (
	"context"
	"gin-mail/dao"
	"gin-mail/pkg/e"
	"gin-mail/pkg/utils"
	"gin-mail/serializer"
)

type CategoryService struct {
}

func (service *CategoryService) List(ctx context.Context) serializer.Response {
	categoryDao := dao.NewCategoryDao(ctx)
	code := e.Success
	categories, err := categoryDao.ListCategory()
	if err != nil {
		utils.LogrusObj.Infoln("Category Service func List: ", err)
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.BuildListResponse(serializer.BuildCategories(categories), uint(len(categories)))
}
