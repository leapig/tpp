package wo

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
	Id() string
	Token() string
	ApiCreatePreAuthCode() (res map[string]interface{})
	ApiQueryAuth(authorizationCode string) (res map[string]interface{})
	ApiGetAuthorizerList() (res []interface{})
	ApiGetAuthorizerInfo(appId string) (res map[string]interface{})
	GetTemplateDraftList() (res map[string]interface{})
	AddToTemplate(draftId, templateType int64) (res map[string]interface{})
	GetTemplateList(templateType int64) (res map[string]interface{})
	DeleteTemplate(templateId int64) (res map[string]interface{})
	ModifyWxaServerDomain(action, WxaServerDomain string, IsModifyPublishedTogether bool) (res map[string]interface{}, err error)
	GetDomainConfirmFile() (res map[string]interface{}, err error)
	ModifyWxaJumpDomain(action, WxaJumpH5Domain string, IsModifyPublishedTogether bool) (res map[string]interface{}, err error)
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
	payload, _ := json.Marshal(map[string]interface{}{
		"component_appid": a.config.AppId,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/api_create_preauthcode?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("ApiCreatePreAuthCode:%+v", js)
			res = js.MustMap()
		} else {
			logger.Errorf("ApiCreatePreAuthCode:%+v", err)
		}
	} else {
		logger.Errorf("ApiCreatePreAuthCode:%+v", err)
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
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/api_query_auth?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("ApiQueryAuth:%+v", js)
			res = js.MustMap()
		} else {
			logger.Errorf("ApiQueryAuth:%+v", err)
		}
	} else {
		logger.Errorf("ApiQueryAuth:%+v", err)
	}
	return
}

const (
	batchSize = 500
)

// ApiGetAuthorizerList POST https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_list?access_token=ACCESS_TOKEN

func (a *app) ApiGetAuthorizerList() (res []interface{}) {
	offset := 0
	for {
		resp := a.apiGetAuthorizerList(offset * batchSize)
		total, _ := resp["total_count"].(json.Number).Int64()
		res = append(res, resp["list"].([]interface{})...)
		if total <= int64((offset+1)*batchSize) {
			break
		}
		offset++
	}
	return
}

func (a *app) apiGetAuthorizerList(offset int) (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"component_appid": a.config.AppId,
		"offset":          offset,
		"count":           batchSize,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/api_get_authorizer_list?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("apiGetAuthorizerList:%+v", js)
			res = js.MustMap()
		} else {
			logger.Errorf("apiGetAuthorizerList:%+v", err)
		}
	} else {
		logger.Errorf("apiGetAuthorizerList:%+v", err)
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
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/api_get_authorizer_info?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("ApiGetAuthorizerInfo:%+v", js)
			res = js.MustMap()
		} else {
			logger.Errorf("ApiGetAuthorizerInfo:%+v", err)
		}
	} else {
		logger.Errorf("ApiGetAuthorizerInfo:%+v", err)
	}
	return
}

// GetTemplateDraftList GET https://api.weixin.qq.com/wxa/gettemplatedraftlist?access_token=ACCESS_TOKEN
func (a *app) GetTemplateDraftList() (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	req, _ := http.NewRequest(http.MethodGet, a.server+"/wxa/gettemplatedraftlist?"+params.Encode(), nil)
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("GetTemplateDraftList:%+v", js)
			res = js.MustMap()
		} else {
			logger.Errorf("GetTemplateDraftList:%+v", err)
		}
	} else {
		logger.Errorf("GetTemplateDraftList:%+v", err)
	}
	return
}

// AddToTemplate POST https://api.weixin.qq.com/wxa/addtotemplate?access_token=ACCESS_TOKEN
func (a *app) AddToTemplate(draftId, templateType int64) (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]int64{
		"draft_id":      draftId,
		"template_type": templateType,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/wxa/addtotemplate?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("AddToTemplate:%+v", js)
			res = js.MustMap()
		} else {
			logger.Errorf("AddToTemplate:%+v", err)
		}
	} else {
		logger.Errorf("AddToTemplate:%+v", err)
	}
	return
}

// GetTemplateList GET https://api.weixin.qq.com/wxa/gettemplatelist?access_token=ACCESS_TOKEN
func (a *app) GetTemplateList(templateType int64) (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]int64{
		"template_type": templateType,
	})
	req, _ := http.NewRequest(http.MethodGet, a.server+"/wxa/gettemplatelist?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("GetTemplateList:%+v", js)
			res = js.MustMap()
		} else {
			logger.Errorf("GetTemplateList:%+v", err)
		}
	} else {
		logger.Errorf("GetTemplateList:%+v", err)
	}
	return
}

// DeleteTemplate POST https://api.weixin.qq.com/wxa/deletetemplate?access_token=ACCESS_TOKEN
func (a *app) DeleteTemplate(templateId int64) (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]int64{
		"template_id": templateId,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/wxa/deletetemplate?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("DeleteTemplate:%+v", js)
			res = js.MustMap()
		} else {
			logger.Errorf("DeleteTemplate:%+v", err)
		}
	} else {
		logger.Errorf("DeleteTemplate:%+v", err)
	}
	return
}

// ModifyWxaServerDomain POST https://api.weixin.qq.com/cgi-bin/component/modify_wxa_server_domain?access_token=ACCESS_TOKEN
func (a *app) ModifyWxaServerDomain(action, WxaServerDomain string, IsModifyPublishedTogether bool) (res map[string]interface{}, err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	var body map[string]interface{}
	if strings.ToLower(action) != "get" {
		body = map[string]interface{}{
			"action":                       action,
			"wxa_server_domain":            WxaServerDomain,
			"is_modify_published_together": IsModifyPublishedTogether,
		}
	} else {
		body = map[string]interface{}{
			"action": action,
		}
	}
	payload, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/modify_wxa_server_domain?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("ModifyWxaServerDomain:%+v", js)
			if js.Get("errcode").MustInt() != 0 {
				err = errors.New(js.Get("errmsg").MustString())
			} else {
				res = js.MustMap()
			}
		} else {
			logger.Errorf("ModifyWxaServerDomain:%+v", err)
		}
	} else {
		logger.Errorf("ModifyWxaServerDomain:%+v", err)
	}
	return
}

// GetDomainConfirmFile POST https://api.weixin.qq.com/cgi-bin/component/get_domain_confirmfile?access_token=ACCESS_TOKEN
func (a *app) GetDomainConfirmFile() (res map[string]interface{}, err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/get_domain_confirmfile?"+params.Encode(), nil)
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("GetDomainConfirmFile:%+v", js)
			if js.Get("errcode").MustInt() != 0 {
				err = errors.New(js.Get("errmsg").MustString())
			} else {
				res = js.MustMap()
			}
		} else {
			logger.Errorf("GetDomainConfirmFile:%+v", err)
		}
	} else {
		logger.Errorf("GetDomainConfirmFile:%+v", err)
	}
	return
}

// ModifyWxaJumpDomain POST https://api.weixin.qq.com/cgi-bin/component/modify_wxa_jump_domain?access_token=ACCESS_TOKEN
func (a *app) ModifyWxaJumpDomain(action, WxaJumpH5Domain string, IsModifyPublishedTogether bool) (res map[string]interface{}, err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	var body map[string]interface{}
	if strings.ToLower(action) != "get" {
		body = map[string]interface{}{
			"action":                       action,
			"wxa_jump_h5_domain":           WxaJumpH5Domain,
			"is_modify_published_together": IsModifyPublishedTogether,
		}
	} else {
		body = map[string]interface{}{
			"action": action,
		}
	}
	payload, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/modify_wxa_jump_domain?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("ModifyWxaJumpDomain:%+v", js)
			if js.Get("errcode").MustInt() != 0 {
				err = errors.New(js.Get("errmsg").MustString())
			} else {
				res = js.MustMap()
			}
		} else {
			logger.Errorf("ModifyWxaJumpDomain:%+v", err)
		}
	} else {
		logger.Errorf("ModifyWxaJumpDomain:%+v", err)
	}
	return
}
