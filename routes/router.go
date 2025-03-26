package routes

import (
	api "gin-mail/api/v1"
	"gin-mail/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Cors())                  // 跨域
	r.StaticFS("/static", http.Dir("static")) // 加载静态文件的路径

	// 新建路由组
	v1 := r.Group("/api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "success")
		})

		// 用户操作
		v1.POST("user/register", api.UserRegisterHandler)
		v1.POST("user/login", api.UserLoginHandler)

		// 轮播图
		v1.GET("carousels", api.ListCarousel)

		// 商品操作
		v1.GET("products", api.ListProduct)      // 获取商品列表 (不需要用户登录认证)
		v1.GET("products/:id", api.ShowProduct)  // 获取商品详细信息
		v1.GET("imgs/:id", api.ListProductImg)   // 获取商品对应的图片信息
		v1.GET("categories", api.ListCategories) // 获取商品分类

		authed := v1.Group("/")      // 需要登陆保护
		authed.Use(middleware.JWT()) // JWT 认证中间件
		{
			// 用户操作
			authed.PUT("user", api.UserUpdate)               // 修改昵称
			authed.POST("avatar", api.UploadAvatar)          // 上传头像
			authed.POST("user/sending-email", api.SendEmail) // 发送邮件
			authed.POST("user/valid-email", api.ValidEmail)  // 验证邮箱

			// 显示金额
			authed.POST("money", api.ShowMoney) // 显示金额

			// 商品操作
			authed.POST("product", api.CreateProduct)  // 创建商品
			authed.POST("products", api.SearchProduct) // 搜索商品

			// 收藏夹操作
			authed.GET("favorites", api.ListFavorite)          // 获取收藏夹内容
			authed.POST("favorites", api.CreateFavorite)       // 创建收藏夹记录
			authed.DELETE("favorites/:id", api.DeleteFavorite) // 删除收藏夹记录

			// 地址操作
			authed.POST("addresses", api.CreateAddress)       // 创建地址
			authed.GET("addresses/:id", api.ShowAddress)      // 获取地址
			authed.GET("addresses", api.ListAddress)          // 获取用户所有地址
			authed.PUT("addresses/:id", api.UpdateAddress)    // 更新地址
			authed.DELETE("addresses/:id", api.DeleteAddress) // 删除地址
		}
	}
	return r
}
