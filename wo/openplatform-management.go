package wo

import (
	"bytes"
	"encoding/json"
	json2 "github.com/bitly/go-simplejson"
	"net/http"
	"net/url"
)

// BindOpenAccount 绑定开放平台账号
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/openplatform-management/bindOpenAccount.html
// req POST https://api.weixin.qq.com/cgi-bin/open/bind?access_token=ACCESS_TOKEN
func (a *app) BindOpenAccount(authorizerAccessToken, openAppid string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	if openAppid == "" {
		openAppid = a.config.AppId
	}
	payload, _ := json.Marshal(map[string]string{
		"open_appid": openAppid,
	})
	return a.doHttp(http.MethodPost, "/cgi-bin/open/bind?"+params.Encode(), bytes.NewReader(payload))
}

// UnbindOpenAccount 解除绑定开放平台账号
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/openplatform-management/unbindOpenAccount.html
// req POST https://api.weixin.qq.com/cgi-bin/open/unbind?access_token=ACCESS_TOKEN
func (a *app) UnbindOpenAccount(authorizerAccessToken, openAppid string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	if openAppid == "" {
		openAppid = a.config.AppId
	}
	payload, _ := json.Marshal(map[string]string{
		"open_appid": openAppid,
	})
	return a.doHttp(http.MethodPost, "/cgi-bin/open/unbind?"+params.Encode(), bytes.NewReader(payload))
}

// GetOpenAccount 获取开放平台账号
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/openplatform-management/getOpenAccount.html
// req POST https://api.weixin.qq.com/cgi-bin/open/get?access_token=ACCESS_TOKEN
func (a *app) GetOpenAccount(authorizerAccessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	return a.doHttp(http.MethodPost, "/cgi-bin/open/get?"+params.Encode(), nil)
}

// CreateOpenAccount 绑定开放平台账号
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/openplatform-management/createOpenAccount.html
// req POST https://api.weixin.qq.com/cgi-bin/open/create?access_token=ACCESS_TOKEN
func (a *app) CreateOpenAccount(authorizerAccessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	return a.doHttp(http.MethodPost, "/cgi-bin/open/create?"+params.Encode(), nil)
}
