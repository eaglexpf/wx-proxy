package service

import (
	"fmt"

	"math"

	"time"

	"github.com/eaglexpf/wx-proxy/model"
	"github.com/eaglexpf/wx-proxy/util/database"
)

type ScopeService struct {
	scopeModel model.WxScopeModel
}

func (this *ScopeService) Get(wx_id int, uuid string) (data model.WxScopeModel, err error) {
	db := database.DB
	err = db.Where("wx_id=? and uuid=? and is_delete=0", wx_id, uuid).Find(&data).Error
	return
}

func (this *ScopeService) GetList(wx_id, page, page_size int) (data map[string]interface{}, err error) {
	db := database.DB
	var count int
	err = db.Table(this.scopeModel.TableName()).Where("wx_id=? and is_delete=0", wx_id).Count(&count).Error
	if err != nil {
		return
	}
	fmt.Println(count)
	list := []model.WxScopeModel{}
	err = db.Table(this.scopeModel.TableName()).Where("wx_id=? and is_delete=0", wx_id).Limit(page_size).Offset((page - 1) * page_size).Find(&list).Error
	if err != nil {
		return
	}

	data = map[string]interface{}{
		"page":      page,
		"page_size": page_size,
		"count":     count,
		"data":      list,
		"all_page":  math.Ceil(float64(page) / float64(page_size)),
	}
	return
}

func (this *ScopeService) Post(insert_data model.WxScopeModel) (model.WxScopeModel, error) {
	db := database.DB
	err := db.Create(&insert_data).Error
	return insert_data, err
}

func (this *ScopeService) Put(id int, update_data map[string]interface{}) (model.WxScopeModel, error) {
	db := database.DB
	data := this.scopeModel
	err := db.Model(&data).Where("id=? and is_delete=0", id).Update(update_data).Error
	if err != nil {
		return data, err
	}
	err = db.Where("id=? and is_delete=0", id).Find(&data).Error
	return data, err
}

func (this *ScopeService) Delete(id int) error {
	db := database.DB
	data := this.scopeModel
	update_data := map[string]interface{}{
		"is_delete": 1,
		"delete_at": time.Now().Unix(),
	}
	err := db.Model(&data).
		Where("id=? and is_delete=0", id).
		Update(update_data).Error
	return err
}
