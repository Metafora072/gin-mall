package v1

import (
	"gin-mail/pkg/utils"
	"gin-mail/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateCart(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var createCartService service.CartService
	if err := c.ShouldBind(&createCartService); err == nil {
		res := createCartService.Create(c.Request.Context(), claim.Id)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("CreateCart: ", err)
	}
}

func ListCart(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var listCartService service.CartService
	if err := c.ShouldBind(&listCartService); err == nil {
		res := listCartService.List(c.Request.Context(), claim.Id)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("ListCart: ", err)
	}
}

func UpdateCart(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var updateCartService service.CartService
	if err := c.ShouldBind(&updateCartService); err == nil {
		res := updateCartService.Update(c.Request.Context(), claim.Id, c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("UpdateCart: ", err)
	}
}

func DeleteCart(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var deleteCartService service.CartService
	if err := c.ShouldBind(&deleteCartService); err == nil {
		res := deleteCartService.Delete(c.Request.Context(), claim.Id, c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("DeleteCart: ", err)
	}
}
