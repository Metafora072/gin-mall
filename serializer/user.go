package serializer

import (
	"gin-mail/conf"
	"gin-mail/model"
)

// UserVO 是构建后的 User 信息传递给前端 view objective 的结构体
type UserVO struct {
	ID        uint   `json:"id"`
	UserName  string `json:"user_name"`
	NickName  string `json:"nick_name"`
	Type      int    `json:"type"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	Avatar    string `json:"avatar"`
	CreatedAt int64  `json:"created_at"`
}

// BuildUser 将 model.User 构建成 UserVO
func BuildUser(user *model.User) *UserVO {
	return &UserVO{
		ID:       user.ID,
		UserName: user.UserName,
		NickName: user.NickName,
		//Type:
		Email:     user.Email,
		Status:    user.Status,
		Avatar:    conf.Host + conf.HttpPort + conf.AvatarPath + user.Avatar,
		CreatedAt: user.CreatedAt.Unix(),
	}
}
