package model

import (
	"time"

	"github.com/eaglexpf/wx-proxy/config"
	"github.com/jinzhu/gorm"
)

type WxConfig struct {
	ID              int    `gorm:"Colcume:id;PRIMARY_KEY;AUTO_INCREMENT;" json:"-"`
	AppID           string `gorm:"Column:app_id;Type:varchar(100);NOT NULL;UNIQUE_INDEX:app_id;" json:"app_id"`
	AppSecret       string `gorm:"Column:app_secret;Type:varchar(255);NOT NULL;" json:"app_secret"`
	AppName         string `gorm:"Column:app_name;Type:varchar(50);NOT NULL;" json:"app_name"`
	Remark          string `gorm:"Column:remark;Type:varchar(255);" json:"remark"`
	Token           string `gorm:"Column:token;Type:varchar(100);" json:"token"`
	CreateAt        int64  `gorm:"Column:create_at;Type:int(11);NOT NULL;DEFAULT:0;" json:"create_at"`
	UpdateAt        int64  `gorm:"Column:update_at;Type:int(11);NOT NULL;DEFAULT:0;" json:"update_at"`
	AccessToken     string `gorm:"Column:access_token;Type:varchar(512); json:"access_token"`
	AccessTokenTime int64  `gorm:"Column:access_token_time;Type:int(11);NOT NULL;DEFAULT:0;" json:"access_token_time"`
	JsTicket        string `gorm:"Column:js_ticket;Type:varchar(512); json:"js_ticket"`
	JsTicketTime    int64  `gorm:"Column:js_ticket_time;Type:int(11);NOT NULL;DEFAULT:0;" json:"js_ticket_time"`
	IsDelete        int    `gorm:"Column:is_delete;Type:tinyint(2);NOT NULL;DEFAULT:0;UNIQUE_INDEX:app_id;" json:"is_delete"`
	DeleteAt        int64  `gorm:"Column:delete_at;Type:int(11);NOT NULL;DEFAULT:0;UNIQUE_INDEX:app_id;" json:"delete_at"`
}

func (WxConfig) TableName() string {
	return config.Setting.Sql.Profix + "wx_config"
}

func (WxConfig) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return nil
}

func (WxConfig) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdateAt", time.Now().Unix())
	return nil
}
