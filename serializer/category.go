package serializer

import "gin-mail/model"

type CategoryVO struct {
	Id           uint   `json:"id"`
	CategoryName string `json:"category_name"`
	CreatedAt    int64  `json:"created_at"`
}

func BuildCategory(item *model.Category) CategoryVO {
	return CategoryVO{
		Id:           item.ID,
		CategoryName: item.CategoryName,
		CreatedAt:    item.CreatedAt.Unix(),
	}
}

func BuildCategories(items []*model.Category) (categories []CategoryVO) {
	for _, item := range items {
		categories = append(categories, BuildCategory(item))
	}
	return
}
