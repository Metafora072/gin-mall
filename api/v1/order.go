package v1

import (
	"gin-mail/pkg/utils"
	"gin-mail/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateOrder(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var createOrderService service.OrderService
	if err := c.ShouldBind(&createOrderService); err == nil {
		res := createOrderService.Create(c.Request.Context(), claim.Id)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("CreateOrder: ", err)
	}
}

func ListOrder(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var listOrderService service.OrderService
	if err := c.ShouldBind(&listOrderService); err == nil {
		res := listOrderService.List(c.Request.Context(), claim.Id)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("ListOrder: ", err)
	}
}

func ShowOrder(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var showOrderService service.OrderService
	if err := c.ShouldBind(&showOrderService); err == nil {
		res := showOrderService.Show(c.Request.Context(), claim.Id, c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("ShowOrder: ", err)
	}
}

func DeleteOrder(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var deleteOrderService service.OrderService
	if err := c.ShouldBind(&deleteOrderService); err == nil {
		res := deleteOrderService.Delete(c.Request.Context(), claim.Id, c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("DeleteOrder: ", err)
	}
}
