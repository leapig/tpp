package mp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	json2 "github.com/bitly/go-simplejson"
	"github.com/faabiosr/cachego/file"
	"github.com/leapig/tpp/logger"
	"github.com/leapig/tpp/util"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type App interface {
	Key() string
	Id() string
	Token() string
	JsCode2Session(jsCode string) map[string]interface{}
	GetWxACodeUnLimit(page, scene string) []byte
	PostWxaBusinessGetUserPhoneNumber(code string) (res map[string]interface{})
	GetVersionInfo() (res map[string]interface{})
	GetPage() (res map[string]interface{})
}

type GetComponentAccessToken func() string

type Config struct {
	Key               string `json:"key"`
	AppId             string `json:"appid"`
	Secret            string `json:"secret"`
	Version           string `json:"version"`
	ComponentAppid    string `json:"component_appid"`
	GetComponentToken GetComponentAccessToken
}

type app struct {
	config Config
	token  util.AccessToken
	server string
}

func NewApp(config Config) App {
	server := "https://api.weixin.qq.com"
	return &app{
		server: server,
		config: config,
		token: util.AccessToken{
			Id:    config.AppId + config.Secret,
			Cache: file.New(os.TempDir()),
			GetRefreshRequestFunc: func() (resp []byte) {
				if strings.HasPrefix(config.Secret, "refreshtoken@@@") {
					params := url.Values{}
					params.Add("component_access_token", config.GetComponentToken())
					payload, _ := json.Marshal(map[string]string{
						"component_appid":          config.ComponentAppid,
						"authorizer_appid":         config.AppId,
						"authorizer_refresh_token": config.Secret,
					})
					req, _ := http.NewRequest(http.MethodPost, server+"/cgi-bin/component/api_authorizer_token?"+params.Encode(), bytes.NewReader(payload))
					response, _ := http.DefaultClient.Do(req)
					resp, _ = io.ReadAll(response.Body)
				} else {
					params := url.Values{}
					params.Add("appid", config.AppId)
					params.Add("secret", config.Secret)
					params.Add("grant_type", "client_credential")
					req, _ := http.NewRequest(http.MethodGet, server+"/cgi-bin/token?"+params.Encode(), nil)
					response, _ := http.DefaultClient.Do(req)
					resp, _ = io.ReadAll(response.Body)
				}
				fmt.Printf("\n\n%s\n\n", string(resp))
				return
			},
		},
	}
}

// Key 获取当前实例ID
func (a *app) Key() string {
	return a.config.Key
}

// Id 获取当前实例ID
func (a *app) Id() string {
	return a.config.AppId
}

// Token 校验是否配置是否正常（返回access_token）
func (a *app) Token() string {
	return a.token.GetAccessToken()
}

// JsCode2Session
// GET https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code
// GET https://api.weixin.qq.com/sns/component/jscode2session?appid=APPID&js_code=JSCODE&grant_type=authorization_code&component_appid=COMPONENT_APPID&component_access_token=COMPONENT_ACCESS_TOKEN
func (a *app) JsCode2Session(jsCode string) (res map[string]interface{}) {

	if strings.HasPrefix(a.config.Secret, "refreshtoken@@@") {
		params := url.Values{}
		params.Add("component_access_token", a.config.GetComponentToken())
		params.Add("appid", a.config.AppId)
		params.Add("grant_type", "client_credential")
		params.Add("component_appid", a.config.ComponentAppid)
		params.Add("js_code", jsCode)
		if response, err := http.Get(a.server + "/sns/component/jscode2session?" + params.Encode()); err == nil {
			if resp, err := io.ReadAll(response.Body); err == nil {
				js, _ := json2.NewJson(resp)
				logger.Debugf("JsCode2Session:%+v", js)
				if js.Get("openid").MustString() != "" {
					res = js.MustMap()
				}
			}
		}
	} else {
		params := url.Values{}
		params.Add("appid", a.config.AppId)
		params.Add("secret", a.config.Secret)
		params.Add("js_code", jsCode)
		params.Add("grant_type", "authorization_code")
		if response, err := http.Get(a.server + "/sns/jscode2session?" + params.Encode()); err == nil {
			if resp, err := io.ReadAll(response.Body); err == nil {
				js, _ := json2.NewJson(resp)
				logger.Debugf("JsCode2Session:%+v", js)
				if js.Get("openid").MustString() != "" {
					res = js.MustMap()
				}
			}
		}
	}
	return
}

// GetWxACodeUnLimit POST https://api.weixin.qq.com/wxa/getwxacodeunlimit
func (a *app) GetWxACodeUnLimit(page, scene string) []byte {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"page":        page,
		"scene":       scene,
		"check_path":  false,
		"env_version": a.config.Version,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/wxa/getwxacodeunlimit?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			return resp
		}
	}
	return nil
}

// PostWxaBusinessGetUserPhoneNumber POST https://api.weixin.qq.com/wxa/business/getuserphonenumber
func (a *app) PostWxaBusinessGetUserPhoneNumber(code string) (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"code": code,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/wxa/business/getuserphonenumber?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("PostWxaBusinessGetUserPhoneNumber:%+v", js)
			if js.Get("errcode").MustInt() != 0 {
				err = errors.New(js.Get("errmsg").MustString())
			}
			res = js.Get("phone_info").MustMap()
		}
	}
	return
}

// GetAccountBasicInfo POST https://api.weixin.qq.com/cgi-bin/account/getaccountbasicinfo?access_token=ACCESS_TOKEN
func (a *app) GetAccountBasicInfo() (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/account/getaccountbasicinfo?"+params.Encode(), nil)
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("GetAccountBasicInfo:%+v", js)
			res = js.MustMap()
		}
	}
	return
}

// GetVersionInfo POST https://api.weixin.qq.com/wxa/getversioninfo?access_token=ACCESS_TOKEN
func (a *app) GetVersionInfo() (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	req, _ := http.NewRequest(http.MethodPost, a.server+"/wxa/getversioninfo?"+params.Encode(), nil)
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("GetAccountBasicInfo:%+v", js)
			res = js.MustMap()
		}
	}
	return
}

// GetPage GET https://api.weixin.qq.com/wxa/get_page?access_token=ACCESS_TOKEN
func (a *app) GetPage() (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/wxa/get_page?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("GetPage:%+v", js)
			if js.Get("errcode").MustInt() == 0 {
				res = js.MustMap()
			}
		}
	}
	return
}
