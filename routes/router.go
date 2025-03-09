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
			authed.PUT("user", api.UserUpdate)      // 修改昵称
			authed.POST("avatar", api.UploadAvatar) // 上传头像
		}
	}
	return r
}
