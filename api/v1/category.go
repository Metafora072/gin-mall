package v1

import (
	"gin-mail/pkg/utils"
	"gin-mail/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListCategories 是处理获取商品分类路由的 controller 函数
func ListCategories(c *gin.Context) {
	var listCategory service.CategoryService
	if err := c.ShouldBind(&listCategory); err == nil {
		res := listCategory.List(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("ListCategories :", err)
	}
}
