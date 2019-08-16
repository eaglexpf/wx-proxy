package proxy

import (
	"github.com/gin-gonic/gin"
)

type ProxyController struct {
}

func (this *ProxyController) Register(router *gin.Engine) {
	r := router.Group("/proxy")
	{
		r.GET("/", this.Get)
	}
}

func (this *ProxyController) Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
	})
	return
}
