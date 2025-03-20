package v1

import (
	"gin-mail/pkg/utils"
	"gin-mail/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserRegisterHandler 是处理用户注册路由的 controller 函数
func UserRegisterHandler(c *gin.Context) {
	var userRegister service.UserService
	if err := c.ShouldBind(&userRegister); err == nil {
		res := userRegister.Register(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}

}

// UserLoginHandler 是处理用户登录路由的 controller 函数
func UserLoginHandler(c *gin.Context) {
	var userLogin service.UserService
	if err := c.ShouldBind(&userLogin); err == nil {
		res := userLogin.Login(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// UserUpdate 是处理用户修改信息路由的 controller 函数，目前仅支持修改用户昵称 NickName
func UserUpdate(c *gin.Context) {
	var userUpdate service.UserService
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&userUpdate); err == nil {
		res := userUpdate.Update(c.Request.Context(), claims.Id)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// UploadAvatar 是处理用户上传头像路由的 controller 函数
func UploadAvatar(c *gin.Context) {
	// 从 HTTP 请求中提取客户端上传的文件
	// file: 类型为 multipart.File，表示上传文件的 内容流，可直接读取文件数据（如 io.ReadAll(file)）。
	// fileHeader: 类型为 *multipart.FileHeader，包含文件的 元数据（如文件名、大小、MIME 类型等），可通过 fileHeader.Filename 获取原始文件名。
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fileSize := fileHeader.Size
	var uploadAvatar service.UserService
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&uploadAvatar); err == nil {
		res := uploadAvatar.Post(c.Request.Context(), claims.Id, file, fileSize)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// SendEmail 是处理用户发送邮件路由的 controller 函数
func SendEmail(c *gin.Context) {
	var sendEmail service.SendEmailService
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&sendEmail); err == nil {
		res := sendEmail.Send(c.Request.Context(), claims.Id)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// ValidEmail
func ValidEmail(c *gin.Context) {
	var validEmail service.ValidEmailService
	if err := c.ShouldBind(&validEmail); err == nil {
		res := validEmail.Valid(c.Request.Context(), c.GetHeader("Authorization"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}
