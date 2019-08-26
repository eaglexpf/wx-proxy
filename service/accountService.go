package service

import (
	//	"fmt"

	"math"

	"time"

	"github.com/eaglexpf/wx-proxy/model"
	"github.com/eaglexpf/wx-proxy/util/database"
)

type AccountService struct {
	AccountModel     model.WxConfig
	AccountListModel []model.WxConfig
}

func (this *AccountService) GetList(page, page_size int) (map[string]interface{}, error) {
	//	this.Model = model.WxConfig{}
	//	data := this.AccountListModel
	db := database.DB
	var count int
	db.Table(this.AccountModel.TableName()).Where("is_delete=0").Count(&count)
	err := db.Where("is_delete=0").Limit(page_size).Offset((page - 1) * page_size).Find(&this.AccountListModel).Error

	return map[string]interface{}{
		"page":      page,
		"page_size": page_size,
		"count":     count,
		"data":      this.AccountListModel,
		"all_page":  math.Ceil(float64(page) / float64(page_size)),
	}, err
}

func (this *AccountService) GetInfoByAppID(app_id string) (model.WxConfig, error) {
	//	data := model.WxConfig{}
	db := database.DB
	data := this.AccountModel
	//	fmt.Println(data, this.AccountModel)
	err := db.Where("app_id=? and is_delete=0", app_id).Find(&data).Error
	return data, err
}

func (this *AccountService) Post(data model.WxConfig) (model.WxConfig, error) {
	db := database.DB
	err := db.Create(&data).Error
	return data, err
}

func (this *AccountService) Put(id int, update_data map[string]interface{}) (model.WxConfig, error) {
	db := database.DB
	data := this.AccountModel
	err := db.Model(&data).
		Where("id=?", id).
		Update(update_data).Error
	if err != nil {
		return data, err
	}
	err = db.Where("id=? and is_delete=0", id).Find(&data).Error
	return data, err
}

func (this *AccountService) Delete(id int) error {
	db := database.DB
	data := this.AccountModel
	update_data := map[string]interface{}{
		"is_delete": 1,
		"delete_at": time.Now().Unix(),
	}
	err := db.Model(&data).
		Where("id=? and is_delete=0", id).
		Update(update_data).Error
	return err
}
