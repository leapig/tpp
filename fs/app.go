package fs

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
	"time"
)

type App interface {
	Id() string
	Test() string
	TenantQuery() (res map[string]interface{})
	Applications() (res map[string]interface{})
	AppVisibility() map[string]interface{}
	AppContactsRangeConfiguration() map[string]interface{}
	DepartmentListSubId(deptIdList []interface{}) []interface{}
	DepartmentsChildren(departmentId string, pageToken string) []interface{}
	DepartmentGet(id string) (res map[string]interface{})
	UsersFindByDepartment(id string, pageToken string) []interface{}
	UserGet(id string) map[string]interface{}
	UserIdGet(id string) map[string]interface{}
	AppAccessTokenInternal() string
	TicketGet() (ticket string)
	AuthorizationCode(code string) map[string]interface{}
	MessageSend(msg Message) error
}

type Config struct {
	AppID     string `json:"appId"`
	AppSecret string `json:"appSecret"`
}

type app struct {
	config Config
	token  util.AccessToken
	server string
}

func NewApp(config Config) App {
	server := "https://open.feishu.cn"
	// 管理token
	return &app{
		server: server,
		config: config,
		token: util.AccessToken{
			Id:    config.AppID + config.AppSecret,
			Cache: file.New(os.TempDir()),
			GetRefreshRequestFunc: func() []byte {
				payload, _ := json.Marshal(map[string]string{
					"app_id":     config.AppID,
					"app_secret": config.AppSecret,
				})
				req, _ := http.NewRequest(http.MethodPost,
					server+"/open-apis/auth/v3/tenant_access_token/internal", bytes.NewReader(payload))
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
	return a.config.AppID
}

// Test 校验是否配置是否正常（返回access_token）
func (a *app) Test() string {
	return a.token.GetAccessToken()
}

// TenantQuery GET https://open.feishu.cn/open-apis/tenant/v2/tenant/query
func (a *app) TenantQuery() (res map[string]interface{}) {
	req, _ := http.NewRequest(http.MethodGet, a.server+"/open-apis/tenant/v2/tenant/query", nil)
	a.token.SetLarkAccessToken(req.Header)
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("code").MustInt() == 0 {
				res = js.Get("data").Get("tenant").MustMap()
			}
		}
	}
	return
}

// Applications GET /open-apis/application/v6/applications/:app_id
func (a *app) Applications() (res map[string]interface{}) {
	params := url.Values{}
	params.Add("lang", "zh_cn")
	req, _ := http.NewRequest(http.MethodGet, a.server+"/open-apis/application/v6/applications/"+a.config.AppID+"?"+params.Encode(), nil)
	a.token.SetLarkAccessToken(req.Header)
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("code").MustInt() == 0 {
				res = js.Get("data").Get("app").MustMap()
			}
		}
	}
	return
}

// AppVisibility https://open.feishu.cn/open-apis/application/v2/app/visibility
func (a *app) AppVisibility() (res map[string]interface{}) {
	params := url.Values{}
	params.Add("app_id", a.config.AppID)
	req, _ := http.NewRequest(http.MethodGet, a.server+"/open-apis/application/v2/app/visibility?"+params.Encode(), nil)
	a.token.SetLarkAccessToken(req.Header)
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("code").MustInt() == 0 {
				res = js.Get("data").MustMap()
			}
		}
	}
	return res
}

// AppContactsRangeConfiguration GET https://open.feishu.cn/open-apis/application/v6/applications/:app_id/contacts_range_configuration
func (a app) AppContactsRangeConfiguration() (res map[string]interface{}) {
	params := url.Values{}
	params.Add("page_size", "100")
	params.Add("department_id_type", "department_id")
	params.Add("user_id_type", "user_id")
	req, _ := http.NewRequest(http.MethodGet, a.server+"/open-apis/application/v6/applications/"+a.config.AppID+"/contacts_range_configuration", nil)
	a.token.SetLarkAccessToken(req.Header)
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("code").MustInt() == 0 {
				res = js.Get("data").Get("contacts_range").MustMap()
			}
		}
	}
	return
}

func (a *app) DepartmentListSubId(deptIdList []interface{}) (res []interface{}) {
	for _, deptId := range deptIdList {
		deptIdList = append(deptIdList, a.DepartmentsChildren(deptId.(string), "")...)
	}
	res = deptIdList
	return
}

// DepartmentsChildren GET https://open.feishu.cn/open-apis/contact/v3/departments/:department_id/children
func (a *app) DepartmentsChildren(departmentId string, pageToken string) (res []interface{}) {
	params := url.Values{}
	if pageToken != "" {
		params.Add("pageToken", pageToken)
	}
	params.Add("department_id_type", "department_id")
	params.Add("page_size", "50")
	params.Add("fetch_child", "true")
	req, _ := http.NewRequest(http.MethodGet, a.server+"/open-apis/contact/v3/departments/"+departmentId+"/children?"+params.Encode(), nil)
	a.token.SetLarkAccessToken(req.Header)
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("code").MustInt() == 0 {
				for _, item := range js.Get("data").Get("items").MustArray() {
					res = append(res, item.(map[string]interface{})["department_id"])
				}
				if js.Get("data").Get("has_more").MustBool() == true {
					pageToken = js.Get("data").Get("page_token").MustString()
					res = append(res, a.DepartmentsChildren(departmentId, pageToken)...)
				}
			}
		}
	}
	return
}

// DepartmentGet GET https://open.feishu.cn/open-apis/contact/v3/departments/:department_id
func (a *app) DepartmentGet(id string) (res map[string]interface{}) {
	params := url.Values{}
	params.Add("department_id_type", "department_id")
	req, _ := http.NewRequest(http.MethodGet, a.server+"/open-apis/contact/v3/departments/"+id+"?"+params.Encode(), nil)
	a.token.SetLarkAccessToken(req.Header)
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("code").MustInt() == 0 {
				res = js.Get("data").Get("department").MustMap()
			}
		}
	}
	return
}

// UsersFindByDepartment GET https://open.feishu.cn/open-apis/contact/v3/users/find_by_department
func (a *app) UsersFindByDepartment(id string, pageToken string) (res []interface{}) {
	params := url.Values{}
	if pageToken != "" {
		params.Add("page_token", pageToken)
	}
	params.Add("department_id_type", "department_id")
	params.Add("department_id", id)
	params.Add("page_size", "50")
	req, _ := http.NewRequest(http.MethodGet, a.server+"/open-apis/contact/v3/users/find_by_department?"+params.Encode(), nil)
	a.token.SetLarkAccessToken(req.Header)
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("code").MustInt() == 0 {
				if js.Get("data").Get("items") != nil {
					res = js.Get("data").Get("items").MustArray()
				}
				if js.Get("data").Get("has_more").MustBool() == true {
					pageToken = js.Get("data").Get("page_token").MustString()
					res = append(res, a.DepartmentsChildren(id, pageToken)...)
				}
			}
		}
	}
	return
}

// UserGet GET https://open.feishu.cn/open-apis/contact/v3/users/:user_id
func (a *app) UserGet(userId string) (res map[string]interface{}) {
	params := url.Values{}
	params.Add("department_id_type", "department_id")
	params.Add("user_id_type", "user_id")
	req, _ := http.NewRequest(http.MethodGet, a.server+"/open-apis/contact/v3/users/"+userId+"?"+params.Encode(), nil)
	a.token.SetLarkAccessToken(req.Header)
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("code").MustInt() == 0 {
				res = js.Get("data").Get("user").MustMap()
			}
		}
	}
	return
}

// UserIdGet GET https://open.feishu.cn/open-apis/contact/v3/users/:user_id
func (a *app) UserIdGet(openId string) (res map[string]interface{}) {
	params := url.Values{}
	params.Add("department_id_type", "department_id")
	params.Add("user_id_type", "open_id")
	req, _ := http.NewRequest(http.MethodGet, a.server+"/open-apis/contact/v3/users/"+openId+"?"+params.Encode(), nil)
	a.token.SetLarkAccessToken(req.Header)
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("code").MustInt() == 0 {
				res = js.Get("data").Get("user").MustMap()
			}
		}
	}
	return
}

// AppAccessTokenInternal POST https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal
func (a *app) AppAccessTokenInternal() string {
	appAccessToken := ""
	appAccessToken, _ = a.token.Cache.Fetch("app_access_token" + a.token.Id)
	if appAccessToken == "" {
		payload, _ := json.Marshal(map[string]string{
			"app_id":     a.config.AppID,
			"app_secret": a.config.AppSecret,
		})
		req, _ := http.NewRequest(http.MethodPost, a.server+"/open-apis/auth/v3/app_access_token/internal", bytes.NewReader(payload))
		if response, err := http.DefaultClient.Do(req); err == nil {
			if resp, err := io.ReadAll(response.Body); err == nil {
				js, _ := json2.NewJson(resp)
				if js.Get("code").MustInt() == 0 {
					appAccessToken = js.Get("app_access_token").MustString()
					d := time.Duration(js.Get("expire").MustInt()) * time.Second
					_ = a.token.Cache.Save("app_access_token:"+a.token.Id, appAccessToken, d)
				}
			}
		}
	}
	return "Bearer " + appAccessToken
}

// TicketGet POST https://open.feishu.cn/open-apis/jssdk/ticket/get
func (a *app) TicketGet() (ticket string) {
	ticket, _ = a.token.Cache.Fetch("ticket" + a.token.Id)
	if ticket == "" {
		req, _ := http.NewRequest(http.MethodPost, a.server+"/open-apis/jssdk/ticket/get", nil)
		req.Header.Set("Authorization", a.AppAccessTokenInternal())
		if response, err := http.DefaultClient.Do(req); err == nil {
			if resp, err := io.ReadAll(response.Body); err == nil {
				js, _ := json2.NewJson(resp)
				if js.Get("code").MustInt() == 0 {
					ticket = js.Get("data").Get("ticket").MustString()
					d := time.Duration(js.Get("data").Get("expire_in").MustInt()) * time.Second
					_ = a.token.Cache.Save("ticket:"+a.token.Id, ticket, d)
				}
			}
		}
	}
	return
}

func (a *app) AuthorizationCode(code string) (res map[string]interface{}) {
	payload, _ := json.Marshal(map[string]interface{}{
		"grant_type": "authorization_code",
		"code":       code,
	})
	req, _ := http.NewRequest(http.MethodPost, a.server+"/open-apis/authen/v1/access_token", bytes.NewReader(payload))
	req.Header.Set("Authorization", a.AppAccessTokenInternal())
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("code").MustInt() == 0 {
				res = js.Get("data").MustMap()
			}
		}
	}
	return
}

type Message struct {
	Type    string      `json:"msg_type"`
	ToUser  string      `json:"receive_id"`
	Content MessageCard `json:"content"`
}

type MessageMsg struct {
	Type    string `json:"msg_type"`
	ToUser  string `json:"receive_id"`
	Content string `json:"content"`
}

type MessageCard struct {
	Title   string `json:"title"`
	Url     string `json:"url"`
	Content string `json:"content"`
}

// MessageSend POST https://open.feishu.cn/open-apis/im/v1/messages
func (a *app) MessageSend(msg Message) (err error) {
	params := url.Values{}
	params.Add("receive_id_type", "user_id")
	reqMsg := MessageMsg{
		ToUser: msg.ToUser,
	}
	if msg.Type == "" {
		reqMsg.Type = "interactive"
	} else {
		reqMsg.Type = msg.Type
	}
	content, _ := json.Marshal(map[string]interface{}{
		"config": map[string]interface{}{
			"wide_screen_mode": true,
		},
		"header": map[string]interface{}{
			"title": map[string]interface{}{
				"tag":     "plain_text",
				"content": msg.Content.Title,
			},
			"template": "turquoise",
		},
		"elements": []interface{}{
			map[string]interface{}{
				"tag":     "markdown",
				"content": msg.Content.Content,
			},
			map[string]interface{}{
				"tag": "hr",
			},
			map[string]interface{}{
				"tag": "div",
				"text": map[string]interface{}{
					"tag":     "lark_md",
					"content": "[点击查看详情](" + msg.Content.Url + ")",
				},
			},
		},
	})
	reqMsg.Content = string(content)
	payload, _ := json.Marshal(reqMsg)
	req, _ := http.NewRequest(http.MethodPost, a.server+"/open-apis/im/v1/messages?"+params.Encode(), bytes.NewReader(payload))
	req.Header.Set("Authorization", a.AppAccessTokenInternal())
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("code").MustInt() != 0 {
				err = errors.New(js.Get("msg").MustString())
			}
		}
	}
	return
}
