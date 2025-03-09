package e

// 定义的状态码
const (
	Success       = 200
	Error         = 500
	InvalidParams = 400

	ErrorExistUser         = 30001 // 用户已存在
	ErrorFailEncryption    = 30002 // 加密失败
	ErrorExistUserNotFound = 30003 // 用户不存在
	ErrorNotCompare        = 30004 // 密码错误
	ErrorAuthToken         = 30005 // token 认证失败
)
