package util

import (
	"fmt"

	"net/http"

	"strings"

	"io"

	"github.com/eaglexpf/wx-proxy/config"
	"github.com/eaglexpf/wx-proxy/controllers/account"
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

	router.Run(fmt.Sprintf(":%s", config.Setting.BindPort))
}

func registerRouter(router *gin.Engine) {
	router.Use(Cors())
	router.Use(MiddlewareProxy())
	router.NoRoute(NotFound())
	router.NoMethod(NotFound())

	router.LoadHTMLGlob("template/*")

	new(account.AccountController).Register(router)
	new(account.AccessTokenController).Register(router)
	new(scope.ScopeController).Register(router)
	new(user.UserController).Register(router)
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

func MiddlewareProxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		proxy_path := strings.Split(c.Request.URL.Path, "/proxy/")
		if len(proxy_path) == 2 {
			path_arr := strings.Split(proxy_path[1], "/")
			if len(path_arr) >= 2 {
				scheme := path_arr[0]
				host := path_arr[1]
				path := path_arr[2:]
				if scheme == "http" || scheme == "https" {
					Proxy(scheme, host, "/"+strings.Join(path, "/"), c)
				}
			}
		}
		c.Next()
	}
}

func Proxy(scheme, host, path string, c *gin.Context) {
	transport := http.DefaultTransport

	outReq := new(http.Request)
	*outReq = *c.Request

	outReq.Host = host
	outReq.URL.Host = host
	outReq.URL.Scheme = scheme
	outReq.URL.Path = path

	res, err := transport.RoundTrip(outReq)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadGateway)
		c.Abort()
		return
	}
	// 回写http头
	for key, value := range res.Header {
		for _, v := range value {
			c.Writer.Header().Add(key, v)
		}
	}

	c.Writer.WriteHeader(res.StatusCode)
	io.Copy(c.Writer, res.Body)
	res.Body.Close()
	c.Abort()
}
