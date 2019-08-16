package service

import (
	"fmt"

	"errors"

	"time"

	"github.com/eaglexpf/wx-proxy/model"
	wxhttp "github.com/eaglexpf/wx-proxy/util/http"
)

type AccessTokenService struct {
	accountModel   model.WxConfig
	accountService AccountService
}

func (this *AccessTokenService) GetAccessToken(app_id string) (response map[string]interface{}, err error) {
	info, err := this.accountService.GetInfoByAppID(app_id)
	if err != nil {
		return
	}
	if info.ID <= 0 {
		err = errors.New("无效的微信公众号")
		return
	}
	if info.AccessTokenTime > time.Now().Unix() {
		response = map[string]interface{}{
			"access_token":      info.AccessToken,
			"access_token_time": info.AccessTokenTime,
			"timestamp":         info.AccessTokenTime - time.Now().Unix(),
		}
		return
	}
	//网络请求获取access_token
	uri := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	uri = fmt.Sprintf(uri, info.AppID, info.AppSecret)
	request_data, err := wxhttp.HttpGet(uri)
	if err != nil {
		return
	}
	if _, ok := request_data["access_token"]; !ok {
		err = errors.New("从微信服务器获取access_token失败")
		return
	}
	access_token_time := time.Now().Unix() + 6000
	update_data := map[string]interface{}{
		"access_token":      request_data["access_token"],
		"access_token_time": access_token_time,
	}
	info, err = this.accountService.Put(info.ID, update_data)
	if err != nil {
		return
	}
	response = map[string]interface{}{
		"access_token":      info.AccessToken,
		"access_token_time": info.AccessTokenTime,
		"timestamp":         info.AccessTokenTime - time.Now().Unix(),
	}
	return
}
