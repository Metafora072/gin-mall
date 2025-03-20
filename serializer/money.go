package serializer

import (
	"gin-mail/model"
	"gin-mail/pkg/utils"
)

// MoneyVO 是构建后的 Money 信息传递给前端 view objective 的结构体
type MoneyVO struct {
	UserId    uint   `json:"user_id" form:"user_id"`
	UserName  string `json:"user_name" form:"user_name"`
	UserMoney string `json:"user_money" form:"user_money"`
}

// BuildMoney 根据用户信息和传入的 key 构建 MoneyVO 结构体
func BuildMoney(item *model.User, key string) MoneyVO {
	utils.Encrypt.SetKey(key)
	return MoneyVO{
		UserId:    item.ID,
		UserName:  item.UserName,
		UserMoney: utils.Encrypt.AesDecoding(item.Money),
	}
}
