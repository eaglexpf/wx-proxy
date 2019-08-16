package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//发起http get请求，返回数据为json，返回json解析出来的数据
func HttpGet(uri string) (map[string]interface{}, error) {
	req, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	return jsonData, err
}
