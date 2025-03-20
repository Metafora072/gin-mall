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
		}
	}
	return r
}
