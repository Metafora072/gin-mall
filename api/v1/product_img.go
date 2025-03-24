package v1

import (
	"gin-mail/pkg/utils"
	"gin-mail/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListProductImg(c *gin.Context) {
	var listProductImg service.ListProductImg
	if err := c.ShouldBind(&listProductImg); err == nil {
		res := listProductImg.List(c.Request.Context(), c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln("ListProductImg ", err)
	}
}
