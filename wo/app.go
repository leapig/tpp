package wo

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
	"os"
)

type App interface {
	Id() string
	Token() string
	/* authorization-management */

	GetAuthorizerList() (js []*json2.Json, err error)
	GetAuthorizerInfo(authorizerAppId string) (js *json2.Json, err error)
	SetAuthorizerOptionInfo(authorizerAccessToken, optionName, optionValue string) (js *json2.Json, err error)
	GetAuthorizerOptionInfo(authorizerAccessToken, optionName string) (js *json2.Json, err error)
	/* authorization-management */
	/* thirdparty-management */

	GetTemplatedRaftList() (js *json2.Json, err error)
	AddToTemplate(draftId, templateType int64) (js *json2.Json, err error)
	GetTemplateList(templateType int64) (js *json2.Json, err error)
	DeleteTemplate(templateId int64) (js *json2.Json, err error)
	ModifyThirdpartyServerDomain(action, WxaServerDomain string, IsModifyPublishedTogether bool) (js *json2.Json, err error)
	GetThirdpartyJumpDomainConfirmFile() (js *json2.Json, err error)
	ModifyThirdpartyJumpDomain(action, WxaJumpH5Domain string, IsModifyPublishedTogether bool) (js *json2.Json, err error)
	/* thirdparty-management */
	/* ticket-token */

	StartPushTicket() (js *json2.Json, err error)
	GetPreAuthCode() (js *json2.Json, err error)
	GetAuthorizerAccessToken(authorizerAppId, authorizerRefreshToken string) (js *json2.Json, err error)
	GetAuthorizerRefreshToken(authorizationCode string) (js *json2.Json, err error)
	GetComponentAccessToken() (js *json2.Json, err error)
	/* ticket-token */

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

func (a *app) do(req *http.Request) (js *json2.Json, err error) {
	if response, err := http.DefaultClient.Do(req); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			js, _ := json2.NewJson(resp)
			if js.Get("errcode") != nil && js.Get("errcode").MustInt() != 0 {
				err = errors.New(js.Get("errmsg").MustString())
			}
		}
	}
	return
}

func (a *app) Id() string {
	return a.config.AppId
}

func (a *app) Token() string {
	return a.token.GetAccessToken()
}
