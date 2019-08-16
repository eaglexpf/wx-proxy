package config

import (
	"flag"
	"fmt"
	"os"

	"log"

	"github.com/go-ini/ini"
)

type Config struct {
	BindAddress string //http监听地址
	BindPort    string //http监听端口

	Sql SqlConfig //数据库配置
}

type SqlConfig struct {
	Type     string //数据库类型
	Host     string //数据库地址【ip:port】
	User     string //用户名
	Password string //密码
	Name     string //数据库名

	Profix      string //表前缀
	LogMode     bool   //启用Logger，显示详细日志
	MaxIdleConn int    //闲置的连接数
	MaxOpenConn int    //最多打开的连接数
}

const (
	DefaultBindAddress = "0.0.0.0"
	DefaultBindPort    = "8080"

	DefaultSqlType        = "mysql"
	DefaultSqlHost        = "127.0.0.1:3306"
	DefaultSqlName        = "wx"
	DefaultSqlUser        = "root"
	DefaultSqlPwd         = ""
	DefaultSqlProfix      = ""
	DefaultSqlLogMode     = false
	DefaultSqlMaxIdleConn = 10
	DefaultSqlMaxOpenConn = 1000
)

var (
	Setting    = &Config{}
	version    bool
	configFile string
)

const (
	AppVersion = "0.1"
)

func Init() {
	//设置配置文件
	flag.StringVar(&configFile, "c", "", "./main -c /path/conf/conf.ini")
	//查询版本号
	flag.BoolVar(&version, "v", false, "./main -v")
	flag.Parse()
	if version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	Setting = &Config{}
	if configFile == "" {
		Setting.initDefault()
		return
	}

	Setting.parse(configFile)
}

func (this *Config) initDefault() {
	this.BindAddress = DefaultBindAddress
	this.BindPort = DefaultBindPort

	this.Sql.Type = DefaultSqlType
	this.Sql.Host = DefaultSqlHost
	this.Sql.Name = DefaultSqlName
	this.Sql.User = DefaultSqlUser
	this.Sql.Password = DefaultSqlPwd
	this.Sql.Profix = DefaultSqlProfix
	this.Sql.LogMode = DefaultSqlLogMode
	this.Sql.MaxIdleConn = DefaultSqlMaxIdleConn
	this.Sql.MaxOpenConn = DefaultSqlMaxOpenConn
}

func (this *Config) parse(path string) {
	cfg, err := ini.Load(path)
	if err != nil {
		log.Fatalf("无法解析配置文件：%s", err.Error())
	}

	serverSec, err := cfg.GetSection("server")
	if err != nil {
		log.Fatalf("配置文件解析错误：%s", err.Error())
	}
	this.BindPort = serverSec.Key("HttpPort").MustString(DefaultBindPort)

	sqlSec, err := cfg.GetSection("database")
	if err != nil {
		log.Fatalf("配置文件解析错误：%s", err.Error())
	}
	this.Sql.Type = sqlSec.Key("Type").MustString(DefaultSqlType)
	this.Sql.Host = sqlSec.Key("Host").MustString(DefaultSqlHost)
	this.Sql.Name = sqlSec.Key("Name").MustString(DefaultSqlName)
	this.Sql.User = sqlSec.Key("User").MustString(DefaultSqlUser)
	this.Sql.Password = sqlSec.Key("Password").MustString(DefaultSqlPwd)
	this.Sql.Profix = sqlSec.Key("Profix").MustString(DefaultSqlProfix)
	this.Sql.LogMode = sqlSec.Key("LogMode").MustBool(DefaultSqlLogMode)
	this.Sql.MaxIdleConn = sqlSec.Key("MaxIdleConn").MustInt(DefaultSqlMaxIdleConn)
	this.Sql.MaxOpenConn = sqlSec.Key("MaxOpenConn").MustInt(DefaultSqlMaxOpenConn)
}
