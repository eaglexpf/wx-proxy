package model

import (
	"time"

	"github.com/eaglexpf/wx-proxy/config"
	"github.com/jinzhu/gorm"
)

type WxScopeModel struct {
	ID        int    `gorm:"Colcume:id;PRIMARY_KEY;AUTO_INCREMENT;" json:"-"`
	WxID      int    `gorm:"Column:wx_id;Type:int(11);NOT NULL;UNIQUE_INDEX:uuid_unique;" json:"-"`
	UUID      string `gorm:"Column:uuid;Type:varchar(50);NOT NULL;UNIQUE_INDEX:uuid_unique;" json:"uuid"` //生成的唯一标识
	Scope     string `gorm:"Column:scope;Type:varchar(50);NOT NULL;" json:"scope"`                        //授权类型
	NotifyUrl string `gorm:"Column:notify_url;Type:varchar(255);NOT NULL;" json:"notify_url"`             //授权后的回调地址
	Data      string `gorm:"Column:data;Type:varchar(512);" json:"data"`                                  //附带参数
	CreateAt  int64  `gorm:"Column:create_at;Type:int(11);NOT NULL;DEFAULT:0;" json:"create_at"`
	UpdateAt  int64  `gorm:"Column:update_at;Type:int(11);NOT NULL;DEFAULT:0;" json:"update_at"`
	IsDelete  int    `gorm:"Column:is_delete;Type:tinyint(2);NOT NULL;DEFAULT:0;UNIQUE_INDEX:uuid_unique;" json:"is_delete"`
	DeleteAt  int64  `gorm:"Column:delete_at;Type:int(11);NOT NULL;DEFAULT:0;UNIQUE_INDEX:uuid_unique;" json:"delete_at"`
}

func (WxScopeModel) TableName() string {
	return config.Setting.Sql.Profix + "wx_scope"
}

func (WxScopeModel) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return nil
}

func (WxScopeModel) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdateAt", time.Now().Unix())
	return nil
}
