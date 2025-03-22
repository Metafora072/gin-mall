package service

import (
	"gin-mail/conf"
	"gin-mail/pkg/utils"
	"io"
	"mime/multipart"
	"os"
	"strconv"
)

// UploadAvatarToLocalStatic 将用户头像保存到本地存储
// 参数：
//   - file: 上传的文件内容（multipart.File 接口）
//   - userId: 用户唯一标识
//   - userName: 用户名（用于构建文件名）
//
// 返回值：
//   - filePath: 相对文件路径（如 "user123/john.jpg"）
//   - err: 错误信息
func UploadAvatarToLocalStatic(file multipart.File, userId uint, userName string) (filePath string, err error) {
	// 将 userId 转换成字符串类型，用于后续路径拼接
	bId := strconv.Itoa(int(userId))
	// 构建基础存储路径（./static/imgs/avatar/user{userId}/）
	basePath := "." + conf.AvatarPath + "user" + bId + "/"
	// 检查 basePath 目录是否存在，不存在则创建
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	// 构建完整文件路径 （./static/imgs/avatar/user{userId}/{userName}.jpg)
	avatarPath := basePath + userName + ".jpg" // TODO: 上传的头像文件不一定是 jpg 格式，可以提取 file 后缀
	// 读取上传文件内容到内存
	content, err := io.ReadAll(file)
	if err != nil {
		utils.LogrusObj.Infoln("err", err)
		return "", err
	}
	// 将文件内容写入本地存储（权限 0666：所有用户可读写)
	err = os.WriteFile(avatarPath, content, 0666)
	if err != nil {
		utils.LogrusObj.Infoln("err", err)
		return
	}
	// 返回相对路径 (user{userId}/{userName}.jpg)
	return "user" + bId + "/" + userName + ".jpg", nil

}

// UploadProductToLocalStatic 将商品图片保存到本地
func UploadProductToLocalStatic(file multipart.File, userId uint, productName string) (filePath string, err error) {
	// 将 userId 转换成字符串类型，用于后续路径拼接
	bId := strconv.Itoa(int(userId))
	// 构建基础存储路径（./static/imgs/product/boss{userId}/）
	basePath := "." + conf.ProductPath + "boss" + bId + "/"
	// 检查 basePath 目录是否存在，不存在则创建
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	// 构建完整文件路径 （./static/imgs/product/boss{userId}/{productName}.jpg)
	productPath := basePath + productName + ".jpg"
	// 读取上传文件内容到内存
	content, err := io.ReadAll(file)
	if err != nil {
		utils.LogrusObj.Infoln("UploadProductToLocalStatic func io.ReadAll err:", err)
		return "", err
	}
	// 将文件内容写入本地存储（权限 0666：所有用户可读写)
	err = os.WriteFile(productPath, content, 0666)
	if err != nil {
		utils.LogrusObj.Infoln("UploadProductToLocalStatic func WriteFile err:", err)
		return
	}
	// 返回相对路径 (boss{userId}/{productName}.jpg)
	return "boss" + bId + "/" + productName + ".jpg", nil

}

// DirExistOrNot 检查本地 fileAddr 文件夹路径是否存在
func DirExistOrNot(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// CreateDir 创建文件夹
func CreateDir(dirName string) bool {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		utils.LogrusObj.Infoln("Create Dir:", err)
		return false
	}
	return true
}
