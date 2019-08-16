package model

import (
	"time"

	"github.com/eaglexpf/wx-proxy/config"
	"github.com/jinzhu/gorm"
)

type WxUserModel struct {
	ID               int    `gorm:"Column:id;PRIMARY_KEY;AUTO_INCREMENT;" json:"-"`
	WxID             int    `gorm:"Column:wx_id;Type:int(11);NOT NULL;UNIQUE_INDEX:openid;" json:"-"`
	Openid           string `gorm:"Column:openid;Type:varchar(100);NOT NULL;UNIQUE_INDEX:openid;" json:"openid"`
	Unionid          string `gorm:"Column:unionid;Type:varchar(100);INDEX:unionid;" json:"unionid"`
	Nickname         string `gorm:"Column:nickname;Type:varchar(100);" json:"nickname"`
	AvatarUrl        string `gorm:"Column:avatar_url;Type:varchar(255);" json:"avatar_url"`
	Sex              int    `gorm:"Column:sex;Type:tinyint(2);" json:"sex"`
	Province         string `gorm:"Column:province;Type:varchar(100);" json:"province"`
	City             string `gorm:"Column:city;Type:varchar(100);" json:"city"`
	Country          string `gorm:"Column:country;Type:varchar(100);" json:"country"`
	Language         string `gorm:"Column:language;Type:varchar(30)" json:"language"`
	CreateAt         int64  `gorm:"Column:create_at;Type:int(11);NOT NULL;DEFAULT:0;" json:"create_at"`
	UpdateAt         int64  `gorm:"Column:update_at;Type:int(11);NOT NULL;DEFAULT:0;" json:"update_at"`
	AccessToken      string `gorm:"Column:access_token;Type:varchar(512);" json:"access_token"`
	AccessTokenTime  int64  `gorm:"Column:access_token_time;Type:int(11);NOT NULL;DEFAULT:0;" json:"access_token_time"`
	RefreshToken     string `gorm:"Column:refresh_token;Type:varchar(512);" json:"refresh_token"`
	RefreshTokenTime int64  `gorm:"Column:refresh_token_time;Type:int(11);NOT NULL;DEFAULT:0;" json:"refresh_token_time"`
}

func (WxUserModel) TableName() string {
	return config.Setting.Sql.Profix + "wx_user"
}

func (WxUserModel) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("create_at", time.Now().Unix())
	return nil
}

func (WxUserModel) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("update_at", time.Now().Unix())
	return nil
}
