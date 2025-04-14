package wo

import (
	"bytes"
	"encoding/json"
	"errors"
	json2 "github.com/bitly/go-simplejson"
	"io"
	"net/http"
	"net/url"
)

// StartPushTicket 开启推送ticket
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/startPushTicket.html
// req POST https://api.weixin.qq.com/cgi-bin/component/api_start_push_ticket
func (a *app) StartPushTicket() (js *json2.Json, err error) {
	payload, _ := json.Marshal(map[string]interface{}{
		"component_appid":  a.config.AppId,
		"component_secret": a.config.AppId,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/api_start_push_ticket", bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode") != nil && js.Get("errcode").MustInt() != 0 {
				err = errors.New(js.Get("errmsg").MustString())
			}
		}
	}
	return
}

// GetPreAuthCode 获取预授权码
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getPreAuthCode.html
// req POST https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?access_token=ACCESS_TOKEN
func (a *app) GetPreAuthCode() (js *json2.Json, err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"component_appid": a.config.AppId,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/api_create_preauthcode?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode") != nil && js.Get("errcode").MustInt() != 0 {
				err = errors.New(js.Get("errmsg").MustString())
			}
		}
	}
	return
}

// GetAuthorizerAccessToken 获取授权账号调用令牌
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getAuthorizerAccessToken.html
// req POST https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token=ACCESS_TOKEN
func (a *app) GetAuthorizerAccessToken(authorizerAppId, authorizerRefreshToken string) (js *json2.Json, err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"component_appid":          a.config.AppId,
		"authorizer_appid":         authorizerAppId,
		"authorizer_refresh_token": authorizerRefreshToken,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/api_authorizer_token?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode") != nil && js.Get("errcode").MustInt() != 0 {
				err = errors.New(js.Get("errmsg").MustString())
			}
		}
	}
	return
}

// GetAuthorizerRefreshToken 获取刷新令牌
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getAuthorizerRefreshToken.html
// req POST https://api.weixin.qq.com/cgi-bin/component/api_query_auth?access_token=ACCESS_TOKEN
func (a *app) GetAuthorizerRefreshToken(authorizationCode string) (js *json2.Json, err error) {
	params := url.Values{}
	params.Add("component_access_token", a.Token())
	payload, _ := json.Marshal(map[string]interface{}{
		"component_appid":    a.config.AppId,
		"authorization_code": authorizationCode,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/api_query_auth?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode") != nil && js.Get("errcode").MustInt() != 0 {
				err = errors.New(js.Get("errmsg").MustString())
			}
		}
	}
	return
}

// GetComponentAccessToken 获取令牌
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getComponentAccessToken.html
// req POST https://api.weixin.qq.com/cgi-bin/component/api_component_token
func (a *app) GetComponentAccessToken() (js *json2.Json, err error) {
	payload, _ := json.Marshal(map[string]string{
		"component_appid":         a.config.AppId,
		"component_appsecret":     a.config.Secret,
		"component_verify_ticket": a.config.Ticket,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/api_component_token", bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode") != nil && js.Get("errcode").MustInt() != 0 {
				err = errors.New(js.Get("errmsg").MustString())
			}
		}
	}
	return
}
