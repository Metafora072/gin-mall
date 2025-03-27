package e

// 定义的状态码
const (
	Success       = 200
	Error         = 500
	InvalidParams = 400

	// user 模块错误
	ErrorExistUser             = 30001 // 用户已存在
	ErrorFailEncryption        = 30002 // 加密失败
	ErrorExistUserNotFound     = 30003 // 用户不存在
	ErrorNotCompare            = 30004 // 密码错误
	ErrorAuthToken             = 30005 // token 认证失败
	ErrorAuthCheckTokenTimeout = 30006 // token 过期
	ErrorUploadFail            = 30007 // 上传失败
	ErrorSendEmail             = 30008 // 发送邮件失败

	// product 模块错误
	ErrorProductImgUpload = 40001 // 上传商品图片失败

	// 收藏夹错误
	ErrorFavoriteExist          = 50001 // 收藏夹已存在
	ErrorFavoriteDeleteNotFound = 50002 // 相应收藏夹不存在

	// address 模块错误
	ErrorAddressNotFound = 60001 // address 记录不存在

	// 购物车 cart 模块错误
	ErrorProductAndBossNotMatch = 70001 // 购物车商品和卖家不匹配
	ErrorCartNotFound           = 70002 // 购物车记录不存在

	// 订单 order 模块错误
	ErrorOrderAddressNotFound = 80001 // 订单地址不存在
	ErrorOrderProductNotFound = 80002 // 订单商品不存在
	ErrorOrderBossNotFound    = 80003 // 订单卖家不存在
	ErrorOrderNotFound        = 80004 // 相应订单不存在

)
