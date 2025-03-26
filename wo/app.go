package wo

import (
	"bytes"
	"encoding/json"
	"fmt"
	json2 "github.com/bitly/go-simplejson"
	"github.com/faabiosr/cachego/file"
	"github.com/leapig/tpp/logger"
	"github.com/leapig/tpp/util"
	"io"
	"net/http"
	"net/url"
	"os"
)

type App interface {
	Id() string
	Token() string
	ApiCreatePreAuthCode() (res map[string]interface{})
	ApiGetAuthorizerInfo(appId string) (res map[string]interface{})
}

type Config struct {
	AppId  string `json:"appid"`
	Secret string `json:"secret"`
	Token  string `json:"token"`
	AesKey string `json:"aes_key"`
	Ticket string `json:"ticket"`
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
				payload, _ := json.Marshal(map[string]string{
					"component_appid":         config.AppId,
					"component_appsecret":     config.Secret,
					"component_verify_ticket": config.Ticket,
				})
				req, _ := http.NewRequest(http.MethodPost, server+"/cgi-bin/component/api_component_token", bytes.NewReader(payload))
				response, _ := http.DefaultClient.Do(req)
				resp, _ = io.ReadAll(response.Body)
				fmt.Printf("\n\n%s\n\n", string(resp))
				return
			},
		},
	}
}

func (a *app) Id() string {
	return a.config.AppId
}

func (a *app) Token() string {
	return a.token.GetAccessToken()
}

// ApiCreatePreAuthCode POST https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?access_token=ACCESS_TOKEN
func (a *app) ApiCreatePreAuthCode() (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/component/api_create_preauthcode?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("ApiCreatePreAuthCode:%+v", js)
			res = js.MustMap()
		}
	}
	return
}

// ApiQueryAuth POST https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token=COMPONENT_ACCESS_TOKEN
func (a *app) ApiQueryAuth(authorizationCode string) (res map[string]interface{}) {
	params := url.Values{}
	params.Add("component_access_token", a.Token())
	payload, _ := json.Marshal(map[string]interface{}{
		"component_appid":    a.config.AppId,
		"authorization_code": authorizationCode,
	})
	if response, err := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/api_query_auth?"+params.Encode(), bytes.NewReader(payload)); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("apiGetAuthorizerList:%+v", js)
			res = js.MustMap()
		}
	}
	return
}

// ApiGetAuthorizerList POST https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_list?access_token=ACCESS_TOKEN
func (a *app) ApiGetAuthorizerList() (res []map[string]interface{}) {
	offset := 0
	for {
		if resp := a.apiGetAuthorizerList(offset * 500); resp["total_count"].(int) >= ((offset + 1) * 500) {
			res = append(res, resp["list"].([]map[string]interface{})...)
			break
		} else {
			res = append(res, resp["list"].([]map[string]interface{})...)
			offset++
		}
	}
	return res
}

func (a *app) apiGetAuthorizerList(offset int) (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"component_appid": a.config.AppId,
		"offset":          offset,
		"count":           500,
	})
	if response, err := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/api_get_authorizer_list?"+params.Encode(), bytes.NewReader(payload)); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("apiGetAuthorizerList:%+v", js)
			res = js.MustMap()
		}
	}
	return
}

// ApiGetAuthorizerInfo POST https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info?access_token=ACCESS_TOKEN
func (a *app) ApiGetAuthorizerInfo(authorizerAppId string) (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]string{
		"component_appid":  a.config.AppId,
		"authorizer_appid": authorizerAppId,
	})
	if response, err := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/api_get_authorizer_info?"+params.Encode(), bytes.NewReader(payload)); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("ApiGetAuthorizerInfo:%+v", js)
			res = js.MustMap()
		}
	}
	return
}
