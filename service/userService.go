package service

import (
	"fmt"

	//	"encoding/json"

	"time"

	"math"

	"github.com/eaglexpf/wx-proxy/model"
	"github.com/eaglexpf/wx-proxy/util/database"
	httpserver "github.com/eaglexpf/wx-proxy/util/http"
)

const (
	uri_access_token  = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	uri_refresh_token = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	uri_userinfo      = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
)

type UserService struct {
	userModel model.WxUserModel
}

func (this *UserService) Get(wx_id int, openid string) (data model.WxUserModel, err error) {
	db := database.DB
	err = db.Where("wx_id = ?", wx_id).Find(&data).Error
	return
}

func (this *UserService) GetList(wx_id, page, page_size int) (data map[string]interface{}, err error) {
	db := database.DB
	var count int
	err = db.Table(this.userModel.TableName()).Where("wx_id=?", wx_id).Count(&count).Error
	if err != nil {
		return
	}
	fmt.Println(count)
	list := []model.WxUserModel{}
	err = db.Table(this.userModel.TableName()).Where("wx_id=?", wx_id).Limit(page_size).Offset((page - 1) * page_size).Find(&list).Error
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

func (this *UserService) Post(insert_data model.WxUserModel) (model.WxUserModel, error) {
	db := database.DB
	err := db.Create(&insert_data).Error
	return insert_data, err
}

func (this *UserService) Put(id int, update_data map[string]interface{}) (model.WxUserModel, error) {
	db := database.DB
	data := this.userModel
	err := db.Model(&data).Where("id=?", id).Update(update_data).Error
	if err != nil {
		return data, err
	}
	err = db.Where("id=?", id).Find(&data).Error
	return data, err
}

//从微信服务器获取获取用户信息--根据code
func (this *UserService) GetInfoByCode(wx_info model.WxConfig, code string) (user model.WxUserModel, err error) {
	//通过code刷新token信息
	var access_token_response map[string]interface{}
	access_token_response, err = this.refreshTokenByCode(wx_info.AppID, wx_info.AppSecret, code)
	if err != nil {
		fmt.Println("80", err)
		return
	}

	openid := access_token_response["openid"].(string)
	access_token := access_token_response["access_token"].(string)
	refresh_token := access_token_response["refresh_token"].(string)
	access_token_time := time.Now().Unix() + 6000
	refresh_token_time := time.Now().Unix() + 29*24*3600

	user_info, err := this.Get(wx_info.ID, openid)
	if err != nil && user_info.ID <= 0 {
		//新建用户
		insert_data := model.WxUserModel{
			WxID:             wx_info.ID,
			Openid:           openid,
			AccessToken:      access_token,
			RefreshToken:     refresh_token,
			AccessTokenTime:  access_token_time,
			RefreshTokenTime: refresh_token_time,
		}
		user_info, err = this.Post(insert_data)
		if err != nil {
			fmt.Println("102", err)
			return
		}
	}
	//修改用户
	update_data := map[string]interface{}{
		"access_token":       access_token,
		"access_token_time":  access_token_time,
		"refresh_token":      refresh_token,
		"refresh_token_time": refresh_token_time,
	}
	//通过access_token获取用户信息
	var userinfo_response map[string]interface{}
	userinfo_response, err = this.refreshUserInfoByAccessToken(user_info.Openid, user_info.AccessToken)
	if err != nil {
		fmt.Println("116", err)
		return
	}
	fmt.Println(userinfo_response["sex"])
	update_data["nickname"] = userinfo_response["nickname"].(string)
	//	update_data["sex"] = userinfo_response["sex"].(int)
	update_data["province"] = userinfo_response["province"].(string)
	update_data["city"] = userinfo_response["city"].(string)
	update_data["country"] = userinfo_response["country"].(string)
	update_data["avatar_url"] = userinfo_response["headimgurl"].(string)

	user, err = this.Put(user_info.ID, update_data)
	fmt.Println("128", err)
	return

}

//从微信服务器获取获取用户信息--根据本地缓存的用户信息[强制刷新用户信息]
func (this *UserService) GetInfoByUser(app_id string, user_info model.WxUserModel) (user model.WxUserModel, err error) {
	now_timestamp := time.Now().Unix()

	if user_info.RefreshTokenTime < now_timestamp {
		err = fmt.Errorf("秘钥已失效，请重新授权")
		return
	}

	update_data := make(map[string]interface{})

	if user_info.AccessTokenTime < now_timestamp {
		//通过refresh_token刷新access_token
		var refresh_response map[string]interface{}
		refresh_response, err = this.refreshAccessTokenByRefreshToken(app_id, user_info.RefreshToken)
		if err != nil {
			return
		}

		user_info.AccessToken = refresh_response["access_token"].(string)
		user_info.AccessTokenTime = time.Now().Unix() + 6000

		update_data["access_token"] = user_info.AccessToken
		update_data["access_token_time"] = user_info.AccessTokenTime
	}
	//通过access_token获取用户信息
	var userinfo_response map[string]interface{}
	userinfo_response, err = this.refreshUserInfoByAccessToken(user_info.Openid, user_info.AccessToken)
	if err != nil {
		return
	}

	update_data["nickname"] = userinfo_response["nickname"].(string)
	update_data["sex"] = userinfo_response["sex"].(int)
	update_data["province"] = userinfo_response["province"].(string)
	update_data["city"] = userinfo_response["city"].(string)
	update_data["country"] = userinfo_response["country"].(string)
	update_data["avatar_url"] = userinfo_response["headimgurl"].(string)

	user, err = this.Put(user_info.ID, update_data)
	return
}

//通过code刷新token
func (this *UserService) refreshTokenByCode(app_id, app_secret, code string) (access_token_response map[string]interface{}, err error) {
	uri := fmt.Sprintf(uri_access_token, app_id, app_secret, code)
	access_token_response, err = httpserver.HttpGet(uri)
	if err != nil {
		return
	}
	if _, ok := access_token_response["errcode"]; ok {
		err = fmt.Errorf("返回的错误信息:%s", access_token_response["errcode"])
		return
	}
	return
}

//根据app_id和refresh_token刷新用户的access_token
func (this *UserService) refreshAccessTokenByRefreshToken(app_id, refresh_token string) (refresh_response map[string]interface{}, err error) {
	//通过refresh_token刷新access_token
	refresh_uri := fmt.Sprintf(uri_refresh_token, app_id, refresh_token)
	refresh_response, err = httpserver.HttpGet(refresh_uri)
	if err != nil {
		return
	}
	if _, ok := refresh_response["errcode"]; ok {
		err = fmt.Errorf("返回的错误信息:%s", refresh_response["errcode"])
		return
	}

	return
}

//通过access_token刷新用户信息
func (this *UserService) refreshUserInfoByAccessToken(openid, access_token string) (userinfo_response map[string]interface{}, err error) {
	//通过access_token获取用户信息
	user_info_uri := fmt.Sprintf(uri_userinfo, access_token, openid)
	userinfo_response, err = httpserver.HttpGet(user_info_uri)
	if err != nil {
		return
	}
	if _, ok := userinfo_response["errcode"]; ok {
		err = fmt.Errorf("返回的错误信息:%s", userinfo_response["errcode"])
		return
	}
	return
}
