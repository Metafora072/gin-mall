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

// ListProduct 是处理获取商品列表路由的 controller 函数
func ListProduct(c *gin.Context) {
	var listProductService service.ProductService
	if err := c.ShouldBind(&listProductService); err == nil {
		res := listProductService.List(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("ListProduct", err)
	}
}

// SearchProduct 是处理搜索商品路由的 controller 函数
func SearchProduct(c *gin.Context) {
	var searchProductService service.ProductService
	if err := c.ShouldBind(&searchProductService); err == nil {
		res := searchProductService.Search(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("SearchProduct", err)
	}
}

// ShowProduct 是处理获取商品详细信息路由的 controller 函数
func ShowProduct(c *gin.Context) {
	var showProductService service.ProductService
	if err := c.ShouldBind(&showProductService); err == nil {
		res := showProductService.Show(c.Request.Context(), c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("ShowProduct", err)
	}
}
