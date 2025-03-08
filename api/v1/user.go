package v1

import (
	"gin-mail/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserRegisterHandler 是处理 /user/register 路由的函数
func UserRegisterHandler(c *gin.Context) {
	var userRegister service.UserService
	if err := c.ShouldBind(&userRegister); err == nil {
		res := userRegister.Register(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}

}
