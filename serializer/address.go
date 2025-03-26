package serializer

import "gin-mail/model"

type AddressVO struct {
	Id        uint   `json:"id"`
	UserId    uint   `json:"user_id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	CreatedAt int64  `json:"created_at"`
}

func BuildAddress(item *model.Address) AddressVO {
	return AddressVO{
		Id:        item.ID,
		UserId:    item.UserID,
		Name:      item.Name,
		Phone:     item.Phone,
		Address:   item.Address,
		CreatedAt: item.CreatedAt.Unix(),
	}

}

func BuildAddresses(items []*model.Address) (addresses []AddressVO) {
	for _, item := range items {
		addresses = append(addresses, BuildAddress(item))
	}
	return
}
