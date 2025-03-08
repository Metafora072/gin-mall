package model

import "gorm.io/gorm"

// Notice 是描述公告的模型
type Notice struct {
	gorm.Model
	Text string `gorm:"type:text"`
}
