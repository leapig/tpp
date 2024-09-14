package wk

import (
	"bytes"
	"encoding/json"
	"fmt"
	json2 "github.com/bitly/go-simplejson"
	"github.com/faabiosr/cachego/file"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"tpp/util"
)

type App interface {
	Id() string
	Test() string
	OrgEduList(cursor int) (res []interface{})
	GetOrgByIds(ids interface{}) (res []interface{})
	GetOrgUsers(id interface{}, cursor int, fetchChild int) (res []interface{})
	GetUserByCardNumber(cardNumbers interface{}) (res []interface{})
	Search(keyword interface{}) (res []interface{})
	AuthorizationCode(wxCode string, appKey string, appSecret string, redirectUri string) (res string)
	GetUserInfoByOauth(accessToken string) (res string)
}

type Config struct {
	AppID     string `json:"appId"`
	AppSecret string `json:"appSecret"`
	AppCode   string `json:"appCode"`
}

type app struct {
	config Config
	token  util.AccessToken
	server string
}

func NewApp(config Config) App {
	server := "https://open.wecard.qq.com"
	// 管理token
	return &app{
		server: server,
		config: config,
		token: util.AccessToken{
			Id:    config.AppID + config.AppSecret,
			Cache: file.New(os.TempDir()),
			GetRefreshRequestFunc: func() []byte {
				payload, _ := json.Marshal(map[string]string{
					"app_key":    config.AppID,
					"app_secret": config.AppSecret,
					"grant_type": "client_credentials",
					"scope":      "base",
					"ocode":      config.AppCode,
				})
				req, _ := http.NewRequest(http.MethodPost,
					server+"/cgi-bin/oauth2/token", bytes.NewReader(payload))
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
	return a.config.AppID
}

// Test 校验是否配置是否正常（返回access_token）
func (a *app) Test() string {
	return a.token.GetAccessToken()
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

// OrgEduList POST https://open.wecard.qq.com/cgi-bin/user/org-edu-list?access_token=access_token
func (a *app) OrgEduList(cursor int) (res []interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"page":      cursor,
		"page_size": 5000,
	})
	if response, err := http.Post(a.server+"/cgi-bin/user/org-edu-list?"+params.Encode(), "application/json;charset=utf-8", bytes.NewReader(payload)); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			respStr, _ := strconv.Unquote(strings.Replace(strconv.Quote(string(resp)), `\\u`, `\u`, -1))
			js, _ := json2.NewJson([]byte(respStr))
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("organization").MustArray()
				if len(res) > 1 {
					cursor = cursor + 1
					res = append(res, a.OrgEduList(cursor)...)
				}
			} else {
				// TODO
			}
		}
	}
	return
}

// GetOrgByIds POST https://open.wecard.qq.com/cgi-bin/org/get-org-by-ids?access_token=access_token
func (a *app) GetOrgByIds(ids interface{}) (res []interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	id, _ := strconv.Atoi(ids.(string))
	payload, _ := json.Marshal(map[string]interface{}{
		"org_ids": []int{id},
	})
	if response, err := http.Post(a.server+"/cgi-bin/org/get-org-by-ids?"+params.Encode(), "application/json;charset=utf-8", bytes.NewReader(payload)); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			respStr, _ := strconv.Unquote(strings.Replace(strconv.Quote(string(resp)), `\\u`, `\u`, -1))
			js, _ := json2.NewJson([]byte(respStr))
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("organization").MustArray()
			}
		} else {
			//
		}
	}
	return
}

// GetOrgUsers POST https://open.wecard.qq.com/cgi-bin/user/get-org-users?access_token=access_token
func (a *app) GetOrgUsers(id interface{}, cursor int, fetchChild int) (res []interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"page":        cursor,
		"page_size":   5000,
		"org_id":      id,
		"fetch_child": fetchChild, //是否递归获取子组织架构下面的成员：1-是；0-否，默认为否
	})
	if response, err := http.Post(a.server+"/cgi-bin/user/get-org-users?"+params.Encode(), "application/json;charset=utf-8", bytes.NewReader(payload)); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			respStr, _ := strconv.Unquote(strings.Replace(strconv.Quote(string(resp)), `\\u`, `\u`, -1))
			js, _ := json2.NewJson([]byte(respStr))
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("userlist").MustArray()
				if len(res) > 0 {
					cursor = cursor + 1
					res = append(res, a.GetOrgUsers(id, cursor, fetchChild)...)
				}
			}
		} else {
			//
		}
	}
	return
}

// GetUserByCardNumber POST https://open.wecard.qq.com/cgi-bin/user/get-user-by-card-numbers?access_token=access_token
func (a *app) GetUserByCardNumber(cardNumbers interface{}) (res []interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"card_numbers": cardNumbers,
	})
	if response, err := http.Post(a.server+"/cgi-bin/user/get-user-by-card-numbers?"+params.Encode(), "application/json;charset=utf-8", bytes.NewReader(payload)); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			respStr, _ := strconv.Unquote(strings.Replace(strconv.Quote(string(resp)), `\\u`, `\u`, -1))
			js, _ := json2.NewJson([]byte(respStr))
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("userlist").MustArray()
			}
		} else {
			//
		}
	}
	return
}

// Search POST https://open.wecard.qq.com/cgi-bin/user/search?access_token=access_token
func (a *app) Search(keyword interface{}) (res []interface{}) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]interface{}{
		"keywords": keyword,
	})
	if response, err := http.Post(a.server+"/cgi-bin/user/search?"+params.Encode(), "application/json;charset=utf-8", bytes.NewReader(payload)); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			respStr, _ := strconv.Unquote(strings.Replace(strconv.Quote(string(resp)), `\\u`, `\u`, -1))
			js, _ := json2.NewJson([]byte(respStr))
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("userlist").MustArray()
			}
		} else {
			//
		}
	}
	return
}

// AuthorizationCode POST https://open.wecard.qq.com/connect/oauth2/token
func (a *app) AuthorizationCode(wxCode string, appKey string, appSecret string, redirectUri string) (res string) {
	payload, _ := json.Marshal(map[string]interface{}{
		"wxcode":       wxCode,
		"app_key":      appKey,
		"app_secret":   appSecret,
		"grant_type":   "authorization_code",
		"redirect_uri": redirectUri,
	})
	if response, err := http.Post(a.server+"/connect/oauth2/token", "application/json;charset=utf-8", bytes.NewReader(payload)); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			respStr, _ := strconv.Unquote(strings.Replace(strconv.Quote(string(resp)), `\\u`, `\u`, -1))
			js, _ := json2.NewJson([]byte(respStr))
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("access_token").MustString()
			}
		} else {
			//
		}
	}
	return
}

// GetUserInfoByOauth POST https://open.wecard.qq.com/connect/oauth/get-user-info
func (a *app) GetUserInfoByOauth(accessToken string) (res string) {
	payload, _ := json.Marshal(map[string]interface{}{
		"access_token": accessToken,
	})
	if response, err := http.Post(a.server+"/connect/oauth/get-user-info", "application/json;charset=utf-8", bytes.NewReader(payload)); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			respStr, _ := strconv.Unquote(strings.Replace(strconv.Quote(string(resp)), `\\u`, `\u`, -1))
			js, _ := json2.NewJson([]byte(respStr))
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("card_number").MustString()
			}
		} else {
			//
		}
	}
	return
}
