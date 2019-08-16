package account

import (
	//	"fmt"

	//	"time"

	"github.com/eaglexpf/wx-proxy/service"
	//	wxhttp "github.com/eaglexpf/wx-proxy/util/http"
	"github.com/gin-gonic/gin"
)

type AccessTokenController struct {
	//	accountService     service.AccountService
	accessTokenService service.AccessTokenService
}

func (this *AccessTokenController) Register(router *gin.Engine) {
	accessTokenGroup := router.Group("/wx/access_token")
	{
		accessTokenGroup.GET("/:app_id", this.Get)
	}
}

func (this *AccessTokenController) Get(c *gin.Context) {
	app_id := c.Param("app_id")
	response, err := this.accessTokenService.GetAccessToken(app_id)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  err.Error(),
			"data": make(map[string]string),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": response,
	})
	return
}
