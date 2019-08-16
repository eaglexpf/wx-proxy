package model

import (
	"fmt"

	"github.com/eaglexpf/wx-proxy/util/database"
)

func Init() {
	db := database.DB
	wxConfigModel := &WxConfig{}
	if !db.HasTable(wxConfigModel.TableName()) {
		db.CreateTable(wxConfigModel)
	}
	wxScopeModel := &WxScopeModel{}
	fmt.Println(wxScopeModel.TableName())
	if !db.HasTable(wxScopeModel.TableName()) {
		db.CreateTable(wxScopeModel)
	}
	wxUserModel := &WxUserModel{}
	if !db.HasTable(wxUserModel.TableName()) {
		db.CreateTable(wxUserModel)
	}
}
