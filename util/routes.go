package util

import (
	"fmt"

	"net/http"

	"strings"

	"io"

	"github.com/eaglexpf/wx-proxy/config"
	"github.com/eaglexpf/wx-proxy/controllers/account"
	"github.com/eaglexpf/wx-proxy/controllers/proxy"
	"github.com/eaglexpf/wx-proxy/controllers/scope"
	"github.com/eaglexpf/wx-proxy/controllers/user"
	"github.com/eaglexpf/wx-proxy/model"
	"github.com/eaglexpf/wx-proxy/util/database"
	"github.com/gin-gonic/gin"
)

func Init() {
	config.Init()
	database.Init()
	model.Init()
	router := gin.Default()
	registerRouter(router)
	//	fmt.Println(config.Setting)
	router.Run(fmt.Sprintf(":%s", config.Setting.BindPort))
}

func registerRouter(router *gin.Engine) {
	router.Use(Cors())
	router.NoRoute(NotFound())
	router.NoMethod(NotFound())

	router.LoadHTMLGlob("template/*")

	new(account.AccountController).Register(router)
	new(account.AccessTokenController).Register(router)
	new(scope.ScopeController).Register(router)
	new(user.UserController).Register(router)
	new(proxy.ProxyController).Register(router)
}

func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "404 Not Found",
			"data": make(map[string]string),
		})
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		path_arr := strings.Split(c.Request.URL.Path, "/proxy/")
		if len(path_arr) == 2 {
			Proxy(path_arr[1], c)
			return
		}
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func Proxy(path string, c *gin.Context) {
	fmt.Println(path)
	transport := http.DefaultTransport

	outReq := new(http.Request)
	*outReq = *c.Request

	outReq.Host = "localhost:9000"
	outReq.URL.Host = "localhost:9000"
	outReq.URL.Scheme = "http"
	outReq.URL.Path = path

	res, err := transport.RoundTrip(outReq)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusBadGateway)
		c.Abort()
		return
	}
	// 回写http头
	for key, value := range res.Header {
		fmt.Println("range1", key, value)
		for _, v := range value {
			c.Writer.Header().Add(key, v)
		}
	}
	fmt.Println(res)

	c.Writer.WriteHeader(res.StatusCode)
	io.Copy(c.Writer, res.Body)
	res.Body.Close()
	c.Abort()
}
