package v1

import (
	"gin-mail/pkg/utils"
	"gin-mail/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateAddress(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var createAddressService service.AddressService
	if err := c.ShouldBind(&createAddressService); err == nil {
		res := createAddressService.Create(c.Request.Context(), claim.Id)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("CreateAddress: ", err)
	}
}

func ShowAddress(c *gin.Context) {
	var showAddressService service.AddressService
	if err := c.ShouldBind(&showAddressService); err == nil {
		res := showAddressService.Show(c.Request.Context(), c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("ShowAddress: ", err)
	}
}

func ListAddress(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var listAddressService service.AddressService
	if err := c.ShouldBind(&listAddressService); err == nil {
		res := listAddressService.List(c.Request.Context(), claim.Id)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("ListAddress: ", err)
	}
}

func UpdateAddress(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var updateAddressService service.AddressService
	if err := c.ShouldBind(&updateAddressService); err == nil {
		res := updateAddressService.Update(c.Request.Context(), claim.Id, c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("UpdateAddress: ", err)
	}
}

func DeleteAddress(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var deleteAddressService service.AddressService
	if err := c.ShouldBind(&deleteAddressService); err == nil {
		res := deleteAddressService.Delete(c.Request.Context(), claim.Id, c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("DeleteAddress: ", err)
	}
}
