package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 是描述用户的模型
type User struct {
	gorm.Model
	UserName       string `gorm:"unique"`
	Email          string
	PasswordDigest string
	NickName       string
	Status         string
	Avatar         string
	Money          string
}

const (
	PasswordCost        = 12       // 密码加密难度
	Active       string = "active" // 激活用户
)

// SetPassword 对传入的 password 密码进行加密，存入 PasswordDigest 字段中
func (user *User) SetPassword(password string) error {
	// bcrypt.GenerateFromPassword 是 Go 中 bcrypt 包提供的函数，用于将密码加密。该函数接收两个参数：
	// []byte(password)：将传入的密码转换为 字节切片 ，这是 bcrypt 函数所要求的格式。
	// PasswordCost：加密的复杂度。PasswordCost 是一个常量，通常表示加密时使用的盐的强度。
	// 该值越大，加密过程越复杂，但也会增加计算时间。常见的值是 10 到 14，数字越大越安全，但也需要更多的计算资源。
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 检验传入的密码 password 是否正确
func (user *User) CheckPassword(password string) bool {
	// bcrypt.CompareHashAndPassword 是 bcrypt 包提供的一个函数，用于比较加密后的密码（哈希值）和明文密码
	// CompareHashAndPassword 会对比这两个值：
	// 如果加密后的密码与明文密码匹配，函数返回 nil。
	// 如果两者不匹配，函数返回一个错误（err）。
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}
