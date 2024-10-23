package mp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	json2 "github.com/bitly/go-simplejson"
	"github.com/faabiosr/cachego/file"
	"github.com/leapig/tpp/util"
	"io"
	"net/http"
	"net/url"
	"os"
)

type App interface {
	Key() string
	Id() string
	Token() string
	JsCode2Session(jsCode string) map[string]interface{}
	GetWxACodeUnLimit(page, scene string) []byte
	PostWxaBusinessGetUserPhoneNumber(code string) (res map[string]interface{})
}

type Config struct {
	Key     string `json:"key"`
	AppId   string `json:"appid"`
	Secret  string `json:"secret"`
	Version string `json:"version"`
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
				params := url.Values{}
				params.Add("appid", config.AppId)
				params.Add("secret", config.Secret)
				params.Add("grant_type", "client_credential")
				req, _ := http.NewRequest(http.MethodGet, server+"/cgi-bin/token?"+params.Encode(), nil)
				response, _ := http.DefaultClient.Do(req)
				resp, _ = io.ReadAll(response.Body)
				fmt.Printf("\n\n%s\n\n", string(resp))
				return
			}},
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

// JsCode2Session GET https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code
func (a *app) JsCode2Session(jsCode string) (res map[string]interface{}) {
	params := url.Values{}
	params.Add("appid", a.config.AppId)
	params.Add("secret", a.config.Secret)
	params.Add("js_code", jsCode)
	params.Add("grant_type", "authorization_code")
	if response, err := http.Get(a.server + "/sns/jscode2session?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("openid").MustString() != "" {
				res = js.MustMap()
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
			if js.Get("errcode").MustInt() != 0 {
				err = errors.New(js.Get("errmsg").MustString())
			}
			return js.Get("phone_info").MustMap()
		}
	}
	return nil
}
