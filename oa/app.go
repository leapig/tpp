package oa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	json2 "github.com/bitly/go-simplejson"
	"github.com/faabiosr/cachego"
	"github.com/faabiosr/cachego/file"
	"github.com/leapig/tpp/logger"
	"github.com/leapig/tpp/util"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const JSAPI = "jsapi"
const WXCARD = "wx_card"

type App interface {
	Key() string
	Id() string
	Token() string
	GetAccountBasicInfo() (res map[string]interface{})
	QrcodeCreate(scene string, limit bool) (res map[string]interface{})
	TemplateGetAllPrivateTemplate() []interface{}
	TemplateApiAddTemplate(templateIdShort int, keywordNameList []string) (templateId string)
	TemplateDelPrivateTemplate(templateId string) (res bool)
	MessageTemplateSend(msg Message) error
	UserGet() (res []interface{})
	UserInfo(openId string) (res map[string]interface{})
	GetCurrentSelfMenuInfo() (res map[string]interface{})
	MenuCreate(button []Button) (err error)
	MenuDelete() (res bool)
	TicketGetTicket(ticketType string) (ticket string)
	AuthorizationCode(code string) (res map[string]interface{})
	CardCodeDecrypt(encryptCode string) (code string)
	OpenGet() (res string)
	OpenBind(openAppid string) (err error)
	OpenUnBind(openAppid string) (err error)
	OpenCreate() (res map[string]interface{})
}

type GetComponentAccessToken func() string

type Config struct {
	Key            string        `json:"key"`
	AppId          string        `json:"appid"`
	Secret         string        `json:"secret"`
	Token          string        `json:"token"`
	AesKey         string        `json:"aes_key"`
	ComponentAppid string        `json:"component_appid"`
	ComponentToken string        `json:"component_token"`
	Cache          cachego.Cache `json:"cache"`
}

type app struct {
	config Config
	token  util.AccessToken
	server string
}

func NewApp(config Config) App {
	server := "https://api.weixin.qq.com"
	if config.Cache == nil {
		config.Cache = file.New(os.TempDir())
	}
	return &app{
		server: server,
		config: config,
		token: util.AccessToken{
			Id:    config.AppId + config.Secret,
			Cache: config.Cache,
			GetRefreshRequestFunc: func() (resp []byte) {
				if strings.HasPrefix(config.Secret, "refreshtoken@@@") {
					params := url.Values{}
					params.Add("component_access_token", config.ComponentToken)
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

// Key 获取当前实例Key
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

// GetAccountBasicInfo GET https://api.weixin.qq.com/cgi-bin/account/getaccountbasicinfo?access_token=ACCESS_TOKEN
func (a *app) GetAccountBasicInfo() (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/account/getaccountbasicinfo?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("GetAccountBasicInfo:%+v", js)
			res = js.MustMap()
		} else {
			logger.Errorf("GetAccountBasicInfo:%+v", err)
		}
	} else {
		logger.Errorf("GetAccountBasicInfo:%+v", err)
	}
	return
}

// QrcodeCreate https://api.weixin.qq.com/cgi-bin/qrcode/create
func (a *app) QrcodeCreate(scene string, limit bool) (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	actionName := "QR_STR_SCENE"
	if limit {
		actionName = "QR_LIMIT_STR_SCENE"
	}
	payload := []byte(`{ "action_name": "` + actionName + `", "action_info": {"scene": {"scene_str": "` + scene + `"}}}`)
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/qrcode/create?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("QrcodeCreate:%+v", js)
			res = js.MustMap()
		} else {
			logger.Errorf("QrcodeCreate:%+v", err)
		}
	} else {
		logger.Errorf("QrcodeCreate:%+v", err)
	}
	return
}

// TemplateGetAllPrivateTemplate GET https://api.weixin.qq.com/cgi-bin/template/get_all_private_template?access_token=ACCESS_TOKEN
func (a *app) TemplateGetAllPrivateTemplate() (res []interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/template/get_all_private_template?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("TemplateGetAllPrivateTemplate:%+v", js)
			res = js.Get("template_list").MustArray()
		} else {
			logger.Errorf("TemplateGetAllPrivateTemplate:%+v", err)
		}
	} else {
		logger.Errorf("TemplateGetAllPrivateTemplate:%+v", err)
	}
	return
}

// TemplateApiAddTemplate POST https://api.weixin.qq.com/cgi-bin/template/api_add_template?access_token=ACCESS_TOKEN
func (a *app) TemplateApiAddTemplate(templateIdShort int, keywordNameList []string) (templateId string) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"template_id_short": templateIdShort,
		"keyword_name_list": keywordNameList,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/template/api_add_template?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("TemplateApiAddTemplate:%+v", js)
			templateId = js.Get("template_id").MustString()
		} else {
			logger.Errorf("TemplateApiAddTemplate:%+v", err)
		}
	} else {
		logger.Errorf("TemplateApiAddTemplate:%+v", err)
	}
	return
}

// TemplateDelPrivateTemplate POST https://api.weixin.qq.com/cgi-bin/template/del_private_template?access_token=ACCESS_TOKEN
func (a *app) TemplateDelPrivateTemplate(templateId string) (res bool) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"template_id": templateId,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/template/del_private_template?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("TemplateDelPrivateTemplate:%+v", js)
			res = js.Get("errcode").MustInt() == 0
		} else {
			logger.Errorf("TemplateDelPrivateTemplate:%+v", err)
		}
	} else {
		logger.Errorf("TemplateDelPrivateTemplate:%+v", err)
	}
	return
}

// Message 微信模板消息结构体
type Message struct {
	Touser      string      `json:"touser"`
	TemplateId  string      `json:"template_id"`
	Url         string      `json:"url,omitempty"`
	Miniprogram Miniprogram `json:"miniprogram,omitempty"`
	Data        interface{} `json:"data"`
}

// Miniprogram 微信小程序配置定义
type Miniprogram struct {
	Appid    string `json:"appid,omitempty"`
	Pagepath string `json:"pagepath,omitempty"`
}

// MessageTemplateSend POST https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=ACCESS_TOKEN
func (a *app) MessageTemplateSend(msg Message) (err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(msg)
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/message/template/send?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("MessageTemplateSend:%+v", js)
			if js.Get("errcode").MustInt() != 0 {
				err = errors.New(js.Get("errmsg").MustString())
			}
		} else {
			logger.Errorf("MessageTemplateSend:%+v", err)
		}
	} else {
		logger.Errorf("MessageTemplateSend:%+v", err)
	}
	return
}

func (a *app) UserGet() (res []interface{}) {
	resp := a.userGet("")
	total, _ := resp["total"].(json.Number).Int64()
	res = append(res, resp["data"].(map[string]interface{})["openid"].([]interface{})...)
	for {
		if int(total) == len(res) {
			break
		}
		resp = a.userGet(resp["next_openid"].(string))
		res = append(res, resp["data"].(map[string]interface{})["openid"].([]interface{})...)
	}
	return
}

// UserGet GET https://api.weixin.qq.com/cgi-bin/user/get?access_token=ACCESS_TOKEN&next_openid=NEXT_OPENID
func (a *app) userGet(nextOpenid string) (res map[string]interface{}) {
	params := url.Values{}
	if nextOpenid != "" {
		params.Add("next_openid", nextOpenid)
	}
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/user/get?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("UserGet:%+v", js)
			res = js.MustMap()
		} else {
			logger.Errorf("UserGet:%+v", err)
		}
	} else {
		logger.Errorf("UserGet:%+v", err)
	}
	return
}

// UserInfo GET https://api.weixin.qq.com/cgi-bin/user/info?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN
func (a *app) UserInfo(openId string) (res map[string]interface{}) {
	params := url.Values{}
	params.Add("openid", openId)
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/user/info?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("UserInfo:%+v", js)
			res = js.MustMap()
		} else {
			logger.Errorf("UserInfo:%+v", err)
		}
	} else {
		logger.Errorf("UserInfo:%+v", err)
	}
	return
}

// GetCurrentSelfMenuInfo GET https://api.weixin.qq.com/cgi-bin/get_current_selfmenu_info?access_token=ACCESS_TOKEN
func (a *app) GetCurrentSelfMenuInfo() (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/get_current_selfmenu_info?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("GetCurrentSelfMenuInfo:%+v", js)
			res = js.MustMap()
		} else {
			logger.Errorf("GetCurrentSelfMenuInfo:%+v", err)
		}
	} else {
		logger.Errorf("GetCurrentSelfMenuInfo:%+v", err)
	}
	return
}

// Button 公众号菜单结构体
type Button struct {
	SubButton []Button `json:"sub_button,omitempty"`
	Type      string   `json:"type,omitempty"`
	Name      string   `json:"name"`
	Key       string   `json:"key,omitempty"`
	Url       string   `json:"url,omitempty"`
	MediaId   string   `json:"media_id,omitempty"`
	Appid     string   `json:"appid,omitempty"`
	Pagepath  string   `json:"pagepath,omitempty"`
	ArticleId string   `json:"article_id,omitempty"`
}

// MenuCreate POST https://api.weixin.qq.com/cgi-bin/menu/create?access_token=ACCESS_TOKEN
func (a *app) MenuCreate(button []Button) error {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	// 特殊json编码，处理json.Marshal()值中& < >符号转义问题
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	_ = jsonEncoder.Encode(map[string][]Button{"button": button})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/menu/create?"+params.Encode(), bf)
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("MenuCreate:%+v", js)
			if js.Get("errcode").MustInt() != 0 {
				return errors.New(js.Get("errmsg").MustString())
			}
		} else {
			logger.Errorf("MenuCreate:%+v", err)
		}
	} else {
		logger.Errorf("MenuCreate:%+v", err)
	}
	return nil
}

// MenuDelete GET https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=ACCESS_TOKEN
func (a *app) MenuDelete() (res bool) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/menu/delete?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("MenuDelete:%+v", js)
			res = js.Get("errcode").MustInt() == 0
		} else {
			logger.Errorf("MenuDelete:%+v", err)
		}
	} else {
		logger.Errorf("MenuDelete:%+v", err)
	}
	return
}

// TicketGetTicket GET https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=ACCESS_TOKEN&type=jsapi
func (a *app) TicketGetTicket(ticketType string) (ticket string) {
	ticket, _ = a.token.Cache.Fetch("ticket" + a.token.Id)
	if ticket == "" {
		params := url.Values{}
		params.Add("type", ticketType)
		params = a.token.ApplyAccessToken(params)
		if response, err := http.Get(a.server + "/cgi-bin/ticket/getticket?" + params.Encode()); err == nil {
			if resp, err := io.ReadAll(response.Body); err == nil {
				js, _ := json2.NewJson(resp)
				logger.Debugf("TicketGetTicket:%+v", js)
				if js.Get("errcode").MustInt() == 0 {
					ticket = js.Get("ticket").MustString()
					d := time.Duration(js.Get("expires_in").MustInt()) * time.Second
					_ = a.token.Cache.Save("ticket:"+a.token.Id, ticket, d)
				}
			} else {
				logger.Errorf("TicketGetTicket:%+v", err)
			}
		} else {
			logger.Errorf("TicketGetTicket:%+v", err)
		}
	}
	return
}

// AuthorizationCode GET https://api.weixin.qq.com/sns/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code
func (a *app) AuthorizationCode(code string) (res map[string]interface{}) {
	params := url.Values{}
	params.Add("appid", a.config.AppId)
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")
	if strings.HasPrefix(a.config.Secret, "refreshtoken@@@") {
		params.Add("component_appid", a.config.ComponentAppid)
		params.Add("component_access_token", a.config.ComponentToken)
		if response, err := http.Get(a.server + "/sns/oauth2/component/access_token?" + params.Encode()); err == nil {
			if resp, err := io.ReadAll(response.Body); err == nil {
				js, _ := json2.NewJson(resp)
				logger.Debugf("AuthorizationCode:%+v", js)
				if js.Get("Openid") != nil {
					res = js.MustMap()
				}
			} else {
				logger.Errorf("AuthorizationCode:%+v", err)
			}
		} else {
			logger.Errorf("AuthorizationCode:%+v", err)
		}
	} else {
		params.Add("secret", a.config.Secret)
		if response, err := http.Get(a.server + "/sns/oauth2/access_token?" + params.Encode()); err == nil {
			if resp, err := io.ReadAll(response.Body); err == nil {
				js, _ := json2.NewJson(resp)
				logger.Debugf("AuthorizationCode:%+v", js)
				if js.Get("Openid") != nil {
					res = js.MustMap()
				}
			} else {
				logger.Errorf("AuthorizationCode:%+v", err)
			}
		} else {
			logger.Errorf("AuthorizationCode:%+v", err)
		}
	}
	return
}

// CardCodeDecrypt POST https://api.weixin.qq.com/card/code/decrypt?access_token=TOKEN
func (a *app) CardCodeDecrypt(encryptCode string) (code string) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"encrypt_code": encryptCode,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/card/code/decrypt?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("CardCodeDecrypt:%+v", js)
			code = js.Get("code").MustString()
		} else {
			logger.Errorf("CardCodeDecrypt:%+v", err)
		}
	} else {
		logger.Errorf("CardCodeDecrypt:%+v", err)
	}
	return
}

// OpenGet POST https://api.weixin.qq.com/cgi-bin/open/get?access_token=ACCESS_TOKEN
func (a *app) OpenGet() (res string) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"appid": a.config.AppId,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/open/get?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("OpenGet:%+v", js)
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("open_appid").MustString()
			}
		} else {
			logger.Errorf("OpenGet:%+v", err)
		}
	} else {
		logger.Errorf("OpenGet:%+v", err)
	}
	return
}

// OpenBind POST https://api.weixin.qq.com/cgi-bin/open/bind?access_token=ACCESS_TOKEN
func (a *app) OpenBind(openAppid string) (err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"appid":      a.config.AppId,
		"open_appid": openAppid,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/open/bind?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("OpenBind:%+v", js)
			if js.Get("errcode").MustInt() != 0 {
				err = errors.New(js.Get("errmsg").MustString())
			}
		} else {
			logger.Errorf("OpenBind:%+v", err)
		}
	} else {
		logger.Errorf("OpenBind:%+v", err)
	}
	return
}

// OpenUnBind POST https://api.weixin.qq.com/cgi-bin/open/unbind?access_token=ACCESS_TOKEN
func (a *app) OpenUnBind(openAppid string) (err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"appid":      a.config.AppId,
		"open_appid": openAppid,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/open/unbind?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("OpenUnBind:%+v", js)
			if js.Get("errcode").MustInt() != 0 {
				err = errors.New(js.Get("errmsg").MustString())
			}
		} else {
			logger.Errorf("OpenUnBind:%+v", err)
		}
	} else {
		logger.Errorf("OpenUnBind:%+v", err)
	}
	return
}

// OpenCreate POST https://api.weixin.qq.com/cgi-bin/open/create?access_token=ACCESS_TOKEN
func (a *app) OpenCreate() (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"appid": a.config.AppId,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/open/create?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			logger.Debugf("OpenCreate:%+v", js)
			res = js.MustMap()
		} else {
			logger.Errorf("OpenCreate:%+v", err)
		}
	} else {
		logger.Errorf("OpenCreate:%+v", err)
	}
	return
}
