package v1

import (
	"gin-mail/pkg/utils"
	"gin-mail/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateProduct 是处理创建商品路由的 controller 函数
func CreateProduct(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var createProductService service.ProductService
	if err := c.ShouldBind(&createProductService); err == nil {
		res := createProductService.Create(c.Request.Context(), claim.Id, files)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("CreateProduct", err)
	}
}
