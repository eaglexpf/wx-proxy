package account

import (
	//	"fmt"

	//	"time"

	"github.com/eaglexpf/wx-proxy/service"
	//	wxhttp "github.com/eaglexpf/wx-proxy/util/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type JsTicketController struct {
	//	accountService     service.AccountService
	jsTicketService service.JsTicketService
}

func (this *JsTicketController) Register(router *gin.Engine) {
	group := router.Group("/wx/js_ticket")
	{
		group.POST("/:app_id", this.Get)
	}
}

func (this *JsTicketController) Get(c *gin.Context) {
	app_id := c.Param("app_id")
	ticket_url := c.PostForm("url")
	if ticket_url == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "url不能为空",
			"data": make(map[string]string),
		})
		return
	}
	_, err := url.Parse(ticket_url)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  err.Error(),
			"data": make(map[string]string),
		})
		return
	}
	response, err := this.jsTicketService.GetJsTicket(app_id, ticket_url)
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
