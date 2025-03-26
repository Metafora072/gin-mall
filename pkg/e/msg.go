package e

var MsgFlags = map[int]string{
	Success:       "ok",
	Error:         "fail",
	InvalidParams: "参数错误",

	// user 模块错误
	ErrorExistUser:             "用户名已存在",
	ErrorFailEncryption:        "密码加密失败",
	ErrorExistUserNotFound:     "用户不存在",
	ErrorNotCompare:            "密码错误",
	ErrorAuthToken:             "token 认证失败",
	ErrorAuthCheckTokenTimeout: "token 过期",
	ErrorUploadFail:            "上传失败",
	ErrorSendEmail:             "发送邮件失败",

	// product 模块错误
	ErrorProductImgUpload: "上传商品图片失败",

	// 收藏夹模块错误
	ErrorFavoriteExist:          "收藏夹已存在",
	ErrorFavoriteDeleteNotFound: "相应收藏夹不存在",

	// address 模块错误
	ErrorAddressNotFound: "相应 address 记录不存在",
}

// GetMsg 获取状态码对应的信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]

	if !ok {
		return MsgFlags[Error]
	}
	return msg
}
