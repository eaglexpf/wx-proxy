package account

import (
	"fmt"

	"github.com/eaglexpf/wx-proxy/model"
	"github.com/eaglexpf/wx-proxy/service"
	"github.com/gin-gonic/gin"
)

type AccountController struct {
	accountModel   model.WxConfig
	accountService service.AccountService
}

func (this *AccountController) Register(router *gin.Engine) {
	accountGroup := router.Group("/wx/account")
	{
		//获取所有的账号列表
		accountGroup.GET("/", this.GetList)
		//根据AppId获取账号信息
		accountGroup.GET("/:app_id", this.Get)
		//添加账号
		accountGroup.POST("/", this.Post)
		//修改账号
		accountGroup.PUT("/:app_id", this.Put)
		//删除账号
		accountGroup.DELETE("/:app_id", this.Delete)
	}
}

func (this *AccountController) Get(c *gin.Context) {
	app_id := c.Param("app_id")
	data, err := this.accountService.GetInfoByAppID(app_id)
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
		"data": data,
	})
	return
}

func (this *AccountController) GetList(c *gin.Context) {
	data, err := this.accountService.GetList(1, 100)
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
		"data": data,
	})
	return
}

func (this *AccountController) Post(c *gin.Context) {
	fmt.Println("post")
	var request_type struct {
		AppName   string `form:"app_name" json:"app_name" binding:"required"`
		AppID     string `form:"app_id" json:"app_name" binding:"required"`
		AppSecret string `form:"app_secret" json:"app_secret" binding:"required"`
		Remark    string `form:"remark" json:"remark"`
		Token     string `form:"token" json:"token"`
		//		Sign      string `form:"sign" json:"sign" binding:"required"`
	}
	err := c.ShouldBind(&request_type)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  err.Error(),
			"data": make(map[string]string),
		})
		return
	}
	data := this.accountModel
	data.AppID = request_type.AppID
	data.AppSecret = request_type.AppSecret
	data.AppName = request_type.AppName
	data.Remark = request_type.Remark
	data.Token = request_type.Token
	info, _ := this.accountService.GetInfoByAppID(data.AppID)
	fmt.Println(info, err)
	if info.ID > 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "数据已存在",
			"data": make(map[string]string),
		})
		return
	}
	request, err := this.accountService.Post(data)
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
		"data": request,
	})
	return
}

func (this *AccountController) Put(c *gin.Context) {
	app_id := c.Param("app_id")
	info, err := this.accountService.GetInfoByAppID(app_id)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": make(map[string]string),
		})
		return
	}
	if info.ID <= 0 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "app_id不存在",
			"data": make(map[string]string),
		})
		return
	}
	update_data := make(map[string]interface{})
	if app_id, ok := c.GetPostForm("app_id"); ok {
		update_data["app_id"] = app_id
	}
	if app_secret, ok := c.GetPostForm("app_secret"); ok {
		update_data["app_secret"] = app_secret
	}
	if app_name, ok := c.GetPostForm("app_name"); ok {
		update_data["app_name"] = app_name
	}
	if remark, ok := c.GetPostForm("remark"); ok {
		update_data["remark"] = remark
	}
	if token, ok := c.GetPostForm("token"); ok {
		update_data["token"] = token
	}

	data, err := this.accountService.Put(info.ID, update_data)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": make(map[string]string),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "scuuess",
		"data": data,
	})
}

func (this *AccountController) Delete(c *gin.Context) {
	app_id := c.Param("app_id")
	info, err := this.accountService.GetInfoByAppID(app_id)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": make(map[string]string),
		})
		return
	}
	if info.ID <= 0 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "app_id不存在",
			"data": make(map[string]string),
		})
		return
	}
	err = this.accountService.Delete(info.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": make(map[string]string),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": make(map[string]string),
	})
}
