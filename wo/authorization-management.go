package wo

import (
	"bytes"
	"encoding/json"
	json2 "github.com/bitly/go-simplejson"
	"net/http"
	"net/url"
)

const (
	batchSize = 500
)

// GetAuthorizerList 拉取已授权的账号信息
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/authorization-management/getAuthorizerList.html
// req POST https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_list?access_token=ACCESS_TOKEN
func (a *app) GetAuthorizerList() (js []*json2.Json, err error) {
	offset := 0
	for {
		if res, err := a.getAuthorizerList(offset * batchSize); err == nil {
			total := res.Get("total_count").MustInt()
			js = append(js, res.Get("list"))
			if total <= ((offset + 1) * batchSize) {
				break
			}
			offset++
		} else {
			break
		}
	}
	return
}

func (a *app) getAuthorizerList(offset int) (js *json2.Json, err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"component_appid": a.config.AppId,
		"offset":          offset,
		"count":           batchSize,
	})
	return a.doHttp(http.MethodPost, "/cgi-bin/component/api_get_authorizer_list?"+params.Encode(), bytes.NewReader(payload))
}

// GetAuthorizerInfo 获取授权账号详情
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/authorization-management/getAuthorizerInfo.html
// req POST https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info?access_token=ACCESS_TOKEN
func (a *app) GetAuthorizerInfo(authorizerAppId string) (js *json2.Json, err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]string{
		"component_appid":  a.config.AppId,
		"authorizer_appid": authorizerAppId,
	})
	return a.doHttp(http.MethodPost, "/cgi-bin/component/api_get_authorizer_info?"+params.Encode(), bytes.NewReader(payload))
}

// SetAuthorizerOptionInfo 设置授权方选项信息
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/authorization-management/setAuthorizerOptionInfo.html
// req POST https://api.weixin.qq.com/cgi-bin/component/set_authorizer_option?access_token=ACCESS_TOKEN
func (a *app) SetAuthorizerOptionInfo(authorizerAccessToken, optionName, optionValue string) (js *json2.Json, err error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	payload, _ := json.Marshal(map[string]string{
		"option_name":  optionName,
		"option_value": optionValue,
	})
	return a.doHttp(http.MethodPost, a.server+"/cgi-bin/component/set_authorizer_option?"+params.Encode(), bytes.NewReader(payload))
}

// GetAuthorizerOptionInfo 获取授权方选项信息
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/authorization-management/getAuthorizerOptionInfo.html
// req POST https://api.weixin.qq.com/cgi-bin/component/get_authorizer_option?access_token=ACCESS_TOKEN
func (a *app) GetAuthorizerOptionInfo(authorizerAccessToken, optionName string) (js *json2.Json, err error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	payload, _ := json.Marshal(map[string]string{
		"option_name": optionName,
	})
	return a.doHttp(http.MethodPost, a.server+"/cgi-bin/component/get_authorizer_option?"+params.Encode(), bytes.NewReader(payload))
}
