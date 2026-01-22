package ww

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	json2 "github.com/bitly/go-simplejson"
	"github.com/faabiosr/cachego/file"
	"github.com/leapig/tpp/util"
)

type App interface {
	Id() string
	Test() string
	AgentGet() (res map[string]interface{})
	DepartmentSimpleList(id string) []interface{}
	DepartmentGet(id string) map[string]interface{}
	UserList(id string) []interface{}
	UserGet(userId string) map[string]interface{}
	GetUserDetail(userTicket string) map[string]interface{}
	GetJsApiTicket() (ticket string)
	GetUserInfo(code string) (res map[string]interface{})
	MessageSend(msg Message) (err error)
}

type Config struct {
	CorpId     string `json:"corpid"`
	CorpSecret string `json:"corpsecret"`
	AgentId    string `json:"agentid"`
}

type app struct {
	config Config
	token  util.AccessToken
	server string
}

func NewApp(config Config) App {
	server := "https://qyapi.weixin.qq.com"
	// 管理token
	return &app{
		server: server,
		config: config,
		token: util.AccessToken{
			Id:    config.CorpId + config.CorpSecret,
			Cache: file.New(os.TempDir()),
			GetRefreshRequestFunc: func() []byte {
				params := url.Values{}
				params.Add("corpid", config.CorpId)
				params.Add("corpsecret", config.CorpSecret)
				req, _ := http.NewRequest(http.MethodGet, server+"/cgi-bin/gettoken?"+params.Encode(), nil)
				response, _ := http.DefaultClient.Do(req)
				resp, _ := io.ReadAll(response.Body)
				fmt.Printf("\n\n%s\n\n", string(resp))
				return resp
			},
		},
	}
}

// Id 获取当前实例ID
func (a *app) Id() string {
	return a.config.CorpId
}

// Test 校验是否配置是否正常（返回access_token）
func (a *app) Test() string {
	return a.token.GetAccessToken()
}

// AgentGet https://qyapi.weixin.qq.com/cgi-bin/agent/get?access_token=ACCESS_TOKEN&agentid=AGENTID
func (a *app) AgentGet() (res map[string]interface{}) {
	params := url.Values{}
	params.Add("agentid", a.config.AgentId)
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/agent/get?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			res = js.MustMap()
		}
	}
	return
}

// DepartmentSimpleList GET https://qyapi.weixin.qq.com/cgi-bin/department/simplelist?access_token=ACCESS_TOKEN&id=ID
func (a *app) DepartmentSimpleList(id string) (res []interface{}) {
	params := url.Values{}
	params.Add("id", id)
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/department/simplelist?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("department_id").MustArray()
			}
		}
	}
	return
}

// DepartmentGet GET https://qyapi.weixin.qq.com/cgi-bin/department/get?access_token=ACCESS_TOKEN&id=ID
func (a *app) DepartmentGet(id string) (res map[string]interface{}) {
	params := url.Values{}
	params.Add("id", id)
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/department/get?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("department").MustMap()
			} else if js.Get("errcode").MustInt() == 60011 {
				// "no privilege to access/modify contact/party/agent"
				res = nil
			}
		}
	}
	return
}

// UserList GET https://qyapi.weixin.qq.com/cgi-bin/user/list?access_token=ACCESS_TOKEN&department_id=DEPARTMENT_ID
func (a *app) UserList(departmentId string) (res []interface{}) {
	params := url.Values{}
	params.Add("department_id", departmentId)
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/user/list?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("userlist").MustArray()
			} else if js.Get("errcode").MustInt() == 60011 {
				// "no privilege to access/modify contact/party/agent"
				res = js.Get("userlist").MustArray()
			} else {
				res = []interface{}{}
			}
		}
	}
	return
}

// UserGet GET https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=ACCESS_TOKEN&userid=USERID
func (a *app) UserGet(userId string) (res map[string]interface{}) {
	params := url.Values{}
	params.Add("userid", userId)
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/user/get?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode").MustInt() == 0 {
				js.Del("errcode")
				js.Del("errmsg")
				res = js.MustMap()
			} else if js.Get("errcode").MustInt() == 60011 {
				// "no privilege to access/modify contact/party/agent"
				res = nil
			}
		}
	}
	return
}

// GetUserDetail POST https://qyapi.weixin.qq.com/cgi-bin/auth/getuserdetail?access_token=ACCESS_TOKEN
func (a *app) GetUserDetail(userTicket string) (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"user_ticket": userTicket,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/auth/getuserdetail?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode").MustInt() == 0 {
				res = js.MustMap()
			}
		}
	}
	return
}

// GetJsApiTicket GET https://qyapi.weixin.qq.com/cgi-bin/get_jsapi_ticket?access_token=ACCESS_TOKEN
func (a *app) GetJsApiTicket() (ticket string) {
	ticket, _ = a.token.Cache.Fetch("ticket" + a.token.Id)
	if ticket == "" {
		params := url.Values{}
		params = a.token.ApplyAccessToken(params)
		if response, err := http.Get(a.server + "/cgi-bin/get_jsapi_ticket?" + params.Encode()); err == nil {
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

// GetUserInfo GET https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=ACCESS_TOKEN&code=CODE
func (a *app) GetUserInfo(code string) (res map[string]interface{}) {
	params := url.Values{}
	params.Add("code", code)
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/user/getuserinfo?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode").MustInt() == 0 {
				res = js.MustMap()
			}
		}
	}
	return
}

// Message 微信模板消息结构体
type Message struct {
	ToUser   string   `json:"touser"`
	MsgType  string   `json:"msgtype"`
	AgentId  string   `json:"agentid"`
	TextCard TextCard `json:"textcard"`
}

type TextCard struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

// MessageSend POST https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=ACCESS_TOKEN
func (a *app) MessageSend(msg Message) (err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	if msg.AgentId == "" {
		msg.AgentId = a.config.AgentId
	}
	if msg.MsgType == "" {
		msg.MsgType = "textcard"
	}
	payload, _ := json.Marshal(msg)
	req, _ := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/message/send?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode").MustInt() != 0 {
				return errors.New(js.Get("errmsg").MustString())
			}
		}
	}
	return
}
