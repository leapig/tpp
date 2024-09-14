package oa

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
	"strings"
	"time"
)

type App interface {
	Id() string
	Test() string
	GetAccountBasicInfo() (res map[string]interface{})
	QrcodeCreate(scene string) (res map[string]interface{})
	TemplateGetAllPrivateTemplate() []interface{}
	TemplateApiAddTemplate(templateIdShort int, keywordNameList []string) (templateId string)
	TemplateDelPrivateTemplate(templateId string) (res bool)
	MessageTemplateSend(msg Message) error
	UserGet(nextOpenid string) (res map[string]interface{})
	UserInfo(openId string) (res map[string]interface{})
	GetCurrentSelfMenuInfo() (res map[string]interface{})
	MenuCreate(button []Button) (err error)
	MenuDelete() (res bool)
	TicketGetTicket() (ticket string)
	AuthorizationCode(code string) (res map[string]interface{})
}

type GetComponentAccessToken func() string

type Config struct {
	RegionRid         int64  `json:"regionRid"`
	AppId             string `json:"appid"`
	Secret            string `json:"secret"`
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
	// 管理token
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

// Id 获取当前实例ID
func (a *app) Id() string {
	return a.config.AppId
}

// Test 校验是否配置是否正常（返回access_token）
func (a *app) Test() string {
	return a.token.GetAccessToken()
}

// GetAccountBasicInfo GET https://api.weixin.qq.com/cgi-bin/account/getaccountbasicinfo?access_token=ACCESS_TOKEN
func (a *app) GetAccountBasicInfo() (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/account/getaccountbasicinfo?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			res = js.MustMap()
		}
	}
	return
}

// QrcodeCreate https://api.weixin.qq.com/cgi-bin/qrcode/create
func (a *app) QrcodeCreate(scene string) (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload := []byte(`{"action_name": "QR_LIMIT_STR_SCENE", "action_info": {"scene": {"scene_str": "` + scene + `"}}}`)
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/qrcode/create?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			res = js.MustMap()
		}
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
			res = js.Get("template_list").MustArray()
		}
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
			templateId = js.Get("template_id").MustString()
		}
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
			res = js.Get("errcode").MustInt() == 0
		}
	}
	return
}

// Message 微信模板消息结构体
type Message struct {
	Touser      string `json:"touser"`
	TemplateId  string `json:"template_id"`
	Url         string `json:"url,omitempty"`
	Miniprogram struct {
		Appid    string `json:"appid,omitempty"`
		Pagepath string `json:"pagepath,omitempty"`
	} `json:"miniprogram,omitempty"`
	Data interface{} `json:"data"`
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
			if js.Get("errcode").MustInt() != 0 {
				err = errors.New(js.Get("errmsg").MustString())
			}
		}
	}
	return
}

// UserGet GET https://api.weixin.qq.com/cgi-bin/user/get?access_token=ACCESS_TOKEN&next_openid=NEXT_OPENID
func (a *app) UserGet(nextOpenid string) (res map[string]interface{}) {
	params := url.Values{}
	if nextOpenid != "" {
		params.Add("next_openid", nextOpenid)
	}
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/user/get?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			res = js.MustMap()
		}
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
			res = js.MustMap()
		}
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
			res = js.MustMap()
		}
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
			if js.Get("errcode").MustInt() != 0 {
				return errors.New(js.Get("errmsg").MustString())
			}
		}
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
			res = js.Get("errcode").MustInt() == 0
		}
	}
	return
}

// TicketGetTicket GET https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=ACCESS_TOKEN&type=jsapi
func (a *app) TicketGetTicket() (ticket string) {
	ticket, _ = a.token.Cache.Fetch("ticket" + a.token.Id)
	if ticket == "" {
		params := url.Values{}
		params.Add("type", "jsapi")
		params = a.token.ApplyAccessToken(params)
		if response, err := http.Get(a.server + "/cgi-bin/ticket/getticket?" + params.Encode()); err == nil {
			if resp, err := io.ReadAll(response.Body); err == nil {
				js, _ := json2.NewJson(resp)
				if js.Get("errcode").MustInt() == 0 {
					ticket = js.Get("ticket").MustString()
					d := time.Duration(js.Get("expires_in").MustInt()) * time.Second
					_ = a.token.Cache.Save("ticket:"+a.token.Id, ticket, d)
				}
			}
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
		params.Add("component_appid", os.Getenv("WeChatAppId"))
		params.Add("component_access_token", a.config.GetComponentToken())
		if response, err := http.Get(a.server + "/sns/oauth2/component/access_token?" + params.Encode()); err == nil {
			if resp, err := io.ReadAll(response.Body); err == nil {
				js, _ := json2.NewJson(resp)
				if js.Get("Openid") != nil {
					res = js.MustMap()
				}
			}
		}
	} else {
		params.Add("secret", a.config.Secret)
		if response, err := http.Get(a.server + "/sns/oauth2/access_token?" + params.Encode()); err == nil {
			if resp, err := io.ReadAll(response.Body); err == nil {
				js, _ := json2.NewJson(resp)
				if js.Get("Openid") != nil {
					res = js.MustMap()
				}
			}
		}
	}
	return
}
