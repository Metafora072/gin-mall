package v1

import (
	"gin-mail/pkg/utils"
	"gin-mail/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListFavorite(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var listFavoriteService service.FavoriteService
	if err := c.ShouldBind(&listFavoriteService); err == nil {
		res := listFavoriteService.List(c.Request.Context(), claim.Id)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("ListFavorite: ", err)
	}
}

func CreateFavorite(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var createFavoriteService service.FavoriteService
	if err := c.ShouldBind(&createFavoriteService); err == nil {
		res := createFavoriteService.Create(c.Request.Context(), claim.Id)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("CreateFavorite: ", err)
	}
}

func DeleteFavorite(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	var deleteFavoriteService service.FavoriteService
	if err := c.ShouldBind(&deleteFavoriteService); err == nil {
		res := deleteFavoriteService.Delete(c.Request.Context(), claim.Id, c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("DeleteFavorite: ", err)
	}
}
