package serializer

import "gin-mail/model"

type CarouselVO struct {
	Id        uint   `json:"id"`
	ImgPath   string `json:"img_path"`
	ProductId uint   `json:"product_id"`
	CreatedAt int64  `json:"created_at"`
}

func BuildCarousel(item *model.Carousel) CarouselVO {
	return CarouselVO{
		Id:        item.ID,
		ImgPath:   item.ImgPath,
		ProductId: item.ProductId,
		CreatedAt: item.CreatedAt.Unix(),
	}
}

func BuildCarousels(items []model.Carousel) (carousels []CarouselVO) {
	for _, item := range items {
		carousels = append(carousels, BuildCarousel(&item))
	}
	return
}
