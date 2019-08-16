package user

import (
	//	"fmt"

	"github.com/eaglexpf/wx-proxy/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	accountService service.AccountService
	userService    service.UserService
}

func (this *UserController) Register(router *gin.Engine) {
	r := router.Group("/wx/user")
	{
		r.GET("/:app_id", this.GetList)
		r.GET("/:app_id/:openid", this.GetByOpenid)
	}
}

func (this *UserController) GetByOpenid(c *gin.Context) {
	app_id := c.Param("app_id")
	openid := c.Param("openid")
	wx_info, err := this.accountService.GetInfoByAppID(app_id)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  err.Error(),
			"data": make(map[string]string),
		})
		return
	}
	useer_info, err := this.userService.Get(wx_info.ID, openid)
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
		"data": useer_info,
	})
	return

}

func (this *UserController) GetList(c *gin.Context) {
	app_id := c.Param("app_id")
	wx_info, err := this.accountService.GetInfoByAppID(app_id)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  err.Error(),
			"data": make(map[string]string),
		})
		return
	}
	user_list, err := this.userService.GetList(wx_info.ID, 1, 100)
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
		"data": user_list,
	})
	return
}
