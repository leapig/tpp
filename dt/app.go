package dt

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
	"strconv"
	"strings"
	"time"
)

type App interface {
	Id() string
	Test() string
	MicroAppAllApps() (res map[string]interface{})
	MicroAppAppsScopes() (res map[string]interface{})
	AuthScopes() (res map[string]interface{})
	DepartmentListSubId(deptIdList []interface{}) (res []interface{})
	DepartmentGet(id interface{}) (res map[string]interface{})
	UserList(id interface{}, cursor interface{}) (res []interface{})
	UserGet(id interface{}) (res map[string]interface{})
	JsApiTickets() (ticket string)
	GetUserInfo(code string) (res map[string]interface{})
	MessageSend(msg Message) (err error)
}

type Config struct {
	CorpId    string `json:"corpId"`
	AppKey    string `json:"appKey"`
	AppSecret string `json:"appSecret"`
	AgentId   int    `json:"agentId"`
}

type app struct {
	config Config
	token  util.AccessToken
	server string
}

func NewApp(config Config) App {
	server := "https://oapi.dingtalk.com"
	// 管理token
	return &app{
		server: server,
		config: config,
		token: util.AccessToken{
			Id:    config.AppKey + config.AppSecret,
			Cache: file.New(os.TempDir()),
			GetRefreshRequestFunc: func() []byte {
				payload, _ := json.Marshal(map[string]string{
					"appKey":    config.AppKey,
					"appSecret": config.AppSecret,
				})
				req, _ := http.NewRequest(http.MethodPost, "https://api.dingtalk.com/v1.0/oauth2/accessToken", bytes.NewReader(payload))
				req.Header.Set("Content-Type", "application/json")
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
	return a.config.AppKey
}

// Test 校验是否配置是否正常（返回access_token）
func (a *app) Test() string {
	return a.token.GetAccessToken()
}

// MicroAppAllApps GET /v1.0/microApp/allApps
func (a *app) MicroAppAllApps() (res map[string]interface{}) {
	params := url.Values{}
	params.Add("x-acs-dingtalk-access-token", a.token.GetAccessToken())
	if response, err := http.Get("https://api.dingtalk.com/v1.0/microApp/allApps?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			appList := js.Get("appList").MustArray()
			for _, app := range appList {
				agentId, _ := strconv.Atoi(string(app.(map[string]interface{})["agentId"].(json.Number)))
				if agentId == a.config.AgentId {
					res = app.(map[string]interface{})
				}
			}
		}
	}
	return
}

// MicroAppAppsScopes GET /v1.0/microApp/apps/{agentId}/scopes
func (a *app) MicroAppAppsScopes() (res map[string]interface{}) {
	params := url.Values{}
	params.Add("x-acs-dingtalk-access-token", a.token.GetAccessToken())
	if response, err := http.Get("https://api.dingtalk.com/v1.0/microApp/apps/" + strconv.Itoa(a.config.AgentId) + "/scopes?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			res = js.Get("result").MustMap()
		}
	}
	return
}

// AuthScopes GET https://oapi.dingtalk.com/auth/scopes
func (a *app) AuthScopes() (res map[string]interface{}) {
	params := url.Values{}
	params.Add("access_token", a.token.GetAccessToken())
	if response, err := http.Get(a.server + "/auth/scopes?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("auth_org_scopes").MustMap()
			} else if js.Get("errcode").MustInt() == 88 && strings.Index(js.Get("errmsg").MustString(), "subcode=90018") > -1 {
				time.Sleep(time.Second)
				res = a.AuthScopes()
			}
		}
	}
	return
}

// DepartmentGet POST https://oapi.dingtalk.com/topapi/v2/department/get?access_token=ACCESS_TOKEN
func (a *app) DepartmentGet(deptId interface{}) (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"dept_id": deptId,
	})
	if response, err := http.Post(a.server+"/topapi/v2/department/get?"+params.Encode(), "application/json;charset=utf-8", bytes.NewReader(payload)); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("result").MustMap()
			} else if js.Get("errcode").MustInt() == 88 && strings.Index(js.Get("errmsg").MustString(), "subcode=90018") > -1 {
				time.Sleep(time.Second)
				res = a.DepartmentGet(deptId)
			}
		}
	}
	return
}

// DepartmentListSubId POST https://oapi.dingtalk.com/topapi/v2/department/listsubid?access_token=ACCESS_TOKEN
func (a *app) DepartmentListSubId(deptIdList []interface{}) (res []interface{}) {
	for _, deptId := range deptIdList {
		params := url.Values{}
		params = a.token.ApplyAccessToken(params)
		payload, _ := json.Marshal(map[string]interface{}{"dept_id": deptId})
		if response, err := http.Post(a.server+"/topapi/v2/department/listsubid?"+params.Encode(), "application/json;charset=utf-8", bytes.NewReader(payload)); err == nil {
			if resp, err := io.ReadAll(response.Body); err == nil {
				js, _ := json2.NewJson(resp)
				if js.Get("errcode").MustInt() == 0 {
					ids := js.Get("result").Get("dept_id_list").MustArray()
					if len(ids) > 0 {
						ids = a.DepartmentListSubId(ids)
						deptIdList = append(deptIdList, ids...)
					}
				} else if js.Get("errcode").MustInt() == 88 && strings.Index(js.Get("errmsg").MustString(), "subcode=90018") > -1 {
					time.Sleep(time.Second)
					ids := a.DepartmentListSubId([]interface{}{deptId})
					deptIdList = append(deptIdList, ids...)
				}
			}
		}
	}
	res = deptIdList
	return
}

// UserList POST https://oapi.dingtalk.com/topapi/v2/user/list?access_token=ACCESS_TOKEN
func (a *app) UserList(id interface{}, cursor interface{}) (res []interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"dept_id": id,
		"cursor":  cursor,
		"size":    100,
	})
	if response, err := http.Post(a.server+"/topapi/v2/user/list?"+params.Encode(), "application/json;charset=utf-8", bytes.NewReader(payload)); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("result").Get("list").MustArray()
				if js.Get("result").Get("has_more").MustBool() == true {
					time.Sleep(time.Second * 50)
					cursor = js.Get("result").Get("next_cursor").MustInt64()
					res = append(res, a.UserList(id, cursor)...)
				}
			} else if js.Get("errcode").MustInt() == 88 && strings.Index(js.Get("errmsg").MustString(), "subcode=90018") > -1 {
				time.Sleep(time.Second)
				res = append(res, a.UserList(id, cursor)...)
			}
		}
	}
	return
}

// UserGet POST https://oapi.dingtalk.com/topapi/v2/user/get?access_token=ACCESS_TOKEN
func (a *app) UserGet(id interface{}) (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"userid":   id,
		"language": "zh_CN"})
	if response, err := http.Post(a.server+"/topapi/v2/user/get?"+params.Encode(), "application/json;charset=utf-8", bytes.NewReader(payload)); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("result").MustMap()
			} else if js.Get("errcode").MustInt() == 88 && strings.Index(js.Get("errmsg").MustString(), "subcode=90018") > -1 {
				time.Sleep(time.Second)
				res = a.UserGet(id)
			}
		}
	}
	return
}

// JsApiTickets POST https://api.dingtalk.com/v1.0/oauth2/jsapiTickets
func (a *app) JsApiTickets() (ticket string) {
	ticket, _ = a.token.Cache.Fetch("ticket" + a.token.Id)
	if ticket == "" {
		req, _ := http.NewRequest(http.MethodPost, "https://api.dingtalk.com/v1.0/oauth2/jsapiTickets", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-acs-dingtalk-access-token", a.token.GetAccessToken())
		if response, err := http.DefaultClient.Do(req); err == nil {
			if resp, err := io.ReadAll(response.Body); err == nil {
				js, _ := json2.NewJson(resp)
				if js.Get("jsapiTicket") != nil {
					ticket = js.Get("jsapiTicket").MustString()
					d := time.Duration(js.Get("expireIn").MustInt()) * time.Second
					_ = a.token.Cache.Save("ticket:"+a.token.Id, ticket, d)
				}
			}
		}
	}
	return
}

// GetUserInfo POST https://oapi.dingtalk.com/topapi/v2/user/getuserinfo
func (a *app) GetUserInfo(code string) (res map[string]interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"code": code,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/topapi/v2/user/getuserinfo?"+params.Encode(), bytes.NewReader(payload))
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("result").MustMap()
			}
		}
	}
	return
}

type Message struct {
	AgentId string     `json:"agent_id"`
	ToUser  string     `json:"userid_list"`
	Msg     MessageMsg `json:"msg"`
}

type MessageMsg struct {
	Type string         `json:"msgtype"`
	Card MessageMsgCard `json:"action_card"`
}

type MessageMsgCard struct {
	Title    string `json:"title"`
	MarkDown string `json:"markdown"`
	Button   string `json:"single_title"`
	Url      string `json:"single_url"`
}

// MessageSend POST https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2?access_token=ACCESS_TOKEN
func (a *app) MessageSend(msg Message) (err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	if msg.AgentId == "" {
		msg.AgentId = strconv.Itoa(a.config.AgentId)
	}
	if msg.Msg.Type == "" {
		msg.Msg.Type = "action_card"
	}
	if msg.Msg.Card.Button == "" {
		msg.Msg.Card.Button = "详情"
	}
	payload, _ := json.Marshal(msg)
	req, _ := http.NewRequest(http.MethodPost, a.server+"/topapi/message/corpconversation/asyncsend_v2?"+params.Encode(), bytes.NewReader(payload))
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
