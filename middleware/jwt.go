package middleware

import (
	"gin-mail/pkg/e"
	"gin-mail/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// JWT 认证中间件，用于在 Gin 框架中验证请求的 Token，并将用户信息注入上下文供后续处理使用
// 检查请求头中的 Authorization 字段是否存在有效的 JWT。
// 若 Token 缺失或无效，返回 JSON 格式的错误信息，并终止请求。
// 若 Token 有效，解析出用户 ID 和用户名，存入 Gin 的上下文（context.Context）中，供后续业务逻辑使用。
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := e.Success
		// c.GetHeader() 是 gin.Context 的方法，用于快速获取请求头中指定键的值。
		// 通过 c.GetHeader("Authorization") 快速获取 JWT Token
		// Authorization 是 HTTP 协议标准定义的请求头字段，用途是携带身份验证信息。
		token := c.GetHeader("Authorization")

		// 若 Authorization 头未提供 Token，返回 HTTP 404 状态码和错误信息。
		if token == "" {
			code = http.StatusNotFound
			c.JSON(http.StatusNotFound, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "缺少token",
			})
			c.Abort()
			return
		}

		// token 存在，则调用 ParseToken 解析 token
		claims, err := utils.ParseToken(token)
		if err != nil { // 解析错误
			code = e.ErrorAuthToken
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "解析token失败",
			})
			c.Abort()
			return
		}

		// token 已被解析完毕，检查 token 是否过期
		if time.Now().Unix() > claims.ExpiresAt.Time.Unix() {
			code = e.ErrorAuthCheckTokenTimeout
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "token过期",
			})
			c.Abort()
			return
		}

		// 注入用户信息到上下文
		// c.Request 是 *gin.Context 对象的一个属性，类型为 *http.Request，表示原始的 HTTP 请求对象。
		// 通过 c.Request，可以访问 HTTP 请求的所有底层细节，包括：
		// 请求方法（GET、POST 等）
		// URL 和路径参数
		// 请求头（Headers）
		// 请求体（Body）
		// 上下文（Context）传递
		// 其他标准库中的 http.Request 功能。
		// c.Request.WithContext() 会创建一个新的 *http.Request 对象，并将新的上下文附加到请求中。
		// 后续处理函数可以通过 c.Request.Context() 获取更新后的上下文。
		// 通过 c.Request.WithContext 更新请求的上下文
		//c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{Id: claims.Id, UserName: claims.UserName}))
		c.Next()
	}
}
