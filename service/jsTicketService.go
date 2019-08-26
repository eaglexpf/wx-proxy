package service

import (
	"fmt"

	"errors"

	"time"

	"github.com/eaglexpf/wx-proxy/model"
	"github.com/eaglexpf/wx-proxy/util/crypto"
	wxhttp "github.com/eaglexpf/wx-proxy/util/http"
	"github.com/satori/go.uuid"
)

type JsTicketService struct {
	accountModel   model.WxConfig
	accountService AccountService
}

func (this *JsTicketService) GetJsTicket(app_id, url string) (response map[string]interface{}, err error) {
	info, err := this.accountService.GetInfoByAppID(app_id)
	if err != nil {
		return
	}
	if info.ID <= 0 {
		err = errors.New("无效的微信公众号")
		return
	}
	signature_fmt := "jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s"
	nonceStr := uuid.Must(uuid.NewV4())
	timestamp := time.Now().Unix()
	if info.JsTicketTime > time.Now().Unix() {
		response = map[string]interface{}{
			"app_id":    app_id,
			"timestamp": timestamp,
			"nonceStr":  nonceStr,
			"url":       url,
			"signature": crypto.Sha1(fmt.Sprintf(signature_fmt, info.JsTicket, nonceStr, timestamp, url)),
		}
		return
	}
	access_token_service := new(AccessTokenService)
	access_token_data, err := access_token_service.GetAccessToken(app_id)
	if err != nil {
		return
	}
	//网络请求获取jsticket
	uri := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
	uri = fmt.Sprintf(uri, access_token_data["access_token"])
	request_data, err := wxhttp.HttpGet(uri)
	if err != nil {
		return
	}
	if _, ok := request_data["ticket"]; !ok {
		err = errors.New("从微信服务器获取jsticket失败")
		return
	}
	js_ticket_time := time.Now().Unix() + 6000
	update_data := map[string]interface{}{
		"js_ticket":      request_data["ticket"],
		"js_ticket_time": js_ticket_time,
	}
	info, err = this.accountService.Put(info.ID, update_data)
	if err != nil {
		return
	}
	response = map[string]interface{}{
		"app_id":    app_id,
		"timestamp": timestamp,
		"nonceStr":  nonceStr,
		"url":       url,
		"signature": crypto.Sha1(fmt.Sprintf(signature_fmt, info.JsTicket, nonceStr, timestamp, url)),
	}
	return
}
