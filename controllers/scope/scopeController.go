package scope

import (
	"fmt"

	"time"

	"strconv"

	"net/http"

	"net/url"

	"github.com/eaglexpf/wx-proxy/model"
	"github.com/eaglexpf/wx-proxy/service"
	"github.com/eaglexpf/wx-proxy/util/crypto"
	"github.com/gin-gonic/gin"
)

const (
	oauthPageUrl = "http://a.frps.8ee4.com/wx/scope/%s/%s/index.html"
)

type ScopeController struct {
	scopeService   service.ScopeService
	accountService service.AccountService
	userService    service.UserService
}

func (this *ScopeController) Register(router *gin.Engine) {
	r := router.Group("/wx/scope")
	{
		r.GET("/:app_id/", this.GetList)
		r.GET("/:app_id/:uuid", this.Get)
		r.POST("/:app_id/", this.Post)
		r.PUT("/:app_id/:uuid", this.Put)
		r.DELETE("/:app_id/:uuid", this.Delete)

		r.GET("/:app_id/:uuid/index.html", this.Page)
	}
}

func (this *ScopeController) Get(c *gin.Context) {
	app_id := c.Param("app_id")
	uuid := c.Param("uuid")
	wx_info, err := this.accountService.GetInfoByAppID(app_id)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  err.Error(),
			"data": make(map[string]string),
		})
		return
	}
	scope_info, err := this.scopeService.Get(wx_info.ID, uuid)
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
		"data": scope_info,
	})
	return
}

func (this *ScopeController) GetList(c *gin.Context) {
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
	scope_list, err := this.scopeService.GetList(wx_info.ID, 1, 100)
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
		"data": scope_list,
	})
	return
}

func (this *ScopeController) Post(c *gin.Context) {
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
	var request_type struct {
		NotifyUrl string `form:"notify_url" json:"notify_url" binding:"required"`
		Scope     string `form:"scope" json:"scope" binding:"required"`
		Data      string `form:"data" json:"data"`
	}
	err = c.ShouldBind(&request_type)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  err.Error(),
			"data": make(map[string]string),
		})
		return
	}
	//	str_time := strconv.FormatInt(time.Now().Unix(), 10)
	fmt.Println(strconv.FormatInt(time.Now().Unix(), 10))
	insert_data := model.WxScopeModel{
		WxID:      wx_info.ID,
		UUID:      crypto.Md5(request_type.NotifyUrl + strconv.FormatInt(time.Now().Unix(), 10)),
		NotifyUrl: request_type.NotifyUrl,
		Scope:     request_type.Scope,
		Data:      request_type.Data,
	}
	response, err := this.scopeService.Post(insert_data)
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

}

func (this *ScopeController) Put(c *gin.Context) {
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
	uuid := c.Param("uuid")
	scope_info, err := this.scopeService.Get(wx_info.ID, uuid)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  err.Error(),
			"data": make(map[string]string),
		})
		return
	}
	update_data := make(map[string]interface{})
	if notify_url, ok := c.GetPostForm("notify_url"); ok {
		update_data["notify_url"] = notify_url
	}
	if data, ok := c.GetPostForm("data"); ok {
		update_data["data"] = data
	}
	response, err := this.scopeService.Put(scope_info.ID, update_data)
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

}

func (this *ScopeController) Delete(c *gin.Context) {
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
	uuid := c.Param("uuid")
	scope_info, err := this.scopeService.Get(wx_info.ID, uuid)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  err.Error(),
			"data": make(map[string]string),
		})
		return
	}
	err = this.scopeService.Delete(scope_info.ID)
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
		"data": make(map[string]string),
	})
}

func (this *ScopeController) Page(c *gin.Context) {
	app_id := c.Param("app_id")
	uuid := c.Param("uuid")
	wx_info, err := this.accountService.GetInfoByAppID(app_id)
	if err != nil {
		c.HTML(200, "scope_error.tmpl", gin.H{
			"msg": err.Error(),
		})
		return
	}
	scope_info, err := this.scopeService.Get(wx_info.ID, uuid)
	if err != nil {
		c.HTML(200, "scope_error.tmpl", gin.H{
			"msg": err.Error(),
		})
		return
	}
	if code, ok := c.GetQuery("code"); ok {
		//根据code获取用户信息
		user, err := this.userService.GetInfoByCode(wx_info, code)
		if err != nil {
			c.HTML(200, "scope_error.tmpl", gin.H{
				"msg": err.Error(),
			})
			return
		}
		//获取用户信息后跳转至授权配置中的指定网页地址+openid
		u, _ := url.Parse(scope_info.NotifyUrl)
		if u.RawQuery == "" {
			u.RawQuery = "openid=" + user.Openid
		} else {
			u.RawQuery += "&openid=" + user.Openid
		}
		c.Redirect(http.StatusMovedPermanently, u.String())
		return
	}
	//没有code跳转至授权页面
	snsapi := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
	fmt.Println("scope_info", scope_info)
	snsapi = fmt.Sprintf(snsapi, wx_info.AppID, fmt.Sprintf(oauthPageUrl, wx_info.AppID, uuid), scope_info.Scope, uuid)

	c.Redirect(http.StatusMovedPermanently, snsapi)
	return
}
