package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"fmt"

	"os"

	"github.com/eaglexpf/wx-proxy/config"
)

type dbConfig struct{}

var DB = &gorm.DB{}

func Init() {
	db := new(dbConfig)
	db.Run()
}

func (this *dbConfig) Run() {
	var err error
	DB, err = gorm.Open(config.Setting.Sql.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Setting.Sql.User,
		config.Setting.Sql.Password,
		config.Setting.Sql.Host,
		config.Setting.Sql.Name))
	if err != nil {
		fmt.Println("sql connect err:", err)
		os.Exit(0)
	}

	DB.SingularTable(true)                 // 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
	DB.LogMode(config.Setting.Sql.LogMode) // 启用Logger，显示详细日志
	DB.DB().SetMaxIdleConns(config.Setting.Sql.MaxIdleConn)
	DB.DB().SetMaxOpenConns(config.Setting.Sql.MaxOpenConn)
}

func (this *dbConfig) Close() {
	DB.Close()
}
