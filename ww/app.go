package ww

import (
	"fmt"
	json2 "github.com/bitly/go-simplejson"
	"github.com/faabiosr/cachego/file"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"tpp/util"
)

type App interface {
	DepartmentSimpleList(id string) []interface{}
	DepartmentGet(id string) map[string]interface{}
	UserSimpleList(id string) []interface{}
}

type Config struct {
	CorpId     string `json:"corpid"`
	CorpSecret string `json:"corpsecret"`
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
				client := http.Client{}
				response, _ := client.Do(req)
				resp, _ := ioutil.ReadAll(response.Body)
				fmt.Printf("\n\n%s\n\n", string(resp))
				return resp
			},
		},
	}
}

// DepartmentSimpleList GET https://qyapi.weixin.qq.com/cgi-bin/department/simplelist?access_token=ACCESS_TOKEN&id=ID
func (a *app) DepartmentSimpleList(id string) (res []interface{}) {
	params := url.Values{}
	params.Add("id", id)
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/department/simplelist?" + params.Encode()); err == nil {
		if resp, err := ioutil.ReadAll(response.Body); err == nil {
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
		if resp, err := ioutil.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("department").MustMap()
			}
		}
	}
	return
}

// UserSimpleList GET https://qyapi.weixin.qq.com/cgi-bin/user/simplelist?access_token=ACCESS_TOKEN&department_id=DEPARTMENT_ID
func (a *app) UserSimpleList(departmentId string) (res []interface{}) {
	params := url.Values{}
	params.Add("department_id", departmentId)
	params = a.token.ApplyAccessToken(params)
	if response, err := http.Get(a.server + "/cgi-bin/user/simplelist?" + params.Encode()); err == nil {
		if resp, err := ioutil.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode").MustInt() == 0 {
				res = js.Get("userlist").MustArray()
			}
		}
	}
	return
}

// ServiceGetPermanentCode https://qyapi.weixin.qq.com/cgi-bin/service/get_permanent_code?suite_access_token=SUITE_ACCESS_TOKEN
func (a *app) ServiceGetPermanentCode() {

}
