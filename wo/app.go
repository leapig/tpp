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
	"os"
	"sync"
	"time"
)

type App interface {
	Id() string
	Token() string
	/* authorization-management */

	GetAuthorizerList() ([]*json2.Json, error)
	GetAuthorizerInfo(authorizerAppId string) (*json2.Json, error)
	SetAuthorizerOptionInfo(authorizerAccessToken, optionName, optionValue string) (*json2.Json, error)
	GetAuthorizerOptionInfo(authorizerAccessToken, optionName string) (*json2.Json, error)
	/* authorization-management */
	/* thirdparty-management */

	GetTemplatedRaftList() (*json2.Json, error)
	AddToTemplate(draftId, templateType int64) (*json2.Json, error)
	GetTemplateList(templateType int64) (*json2.Json, error)
	DeleteTemplate(templateId int64) (*json2.Json, error)
	ModifyThirdpartyServerDomain(action, WxaServerDomain string, IsModifyPublishedTogether bool) (*json2.Json, error)
	GetThirdpartyJumpDomainConfirmFile() (js *json2.Json, err error)
	ModifyThirdpartyJumpDomain(action, WxaJumpH5Domain string, IsModifyPublishedTogether bool) (*json2.Json, error)
	/* thirdparty-management */
	/* ticket-token */

	StartPushTicket() (*json2.Json, error)
	GetPreAuthCode() (*json2.Json, error)
	GetAuthorizerAccessToken(authorizerAppId, authorizerRefreshToken string) (*json2.Json, error)
	GetAuthorizerRefreshToken(authorizationCode string) (*json2.Json, error)
	GetComponentAccessToken() (*json2.Json, error)
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

var (
	httpClient     *http.Client
	httpClientOnce sync.Once
)

func createHTTPClient() *http.Client {
	httpClientOnce.Do(func() {
		// 自定义 Transport 配置连接池
		transport := &http.Transport{
			MaxIdleConns:          100,              // 最大空闲连接数
			MaxIdleConnsPerHost:   10,               // 每个主机的最大空闲连接
			IdleConnTimeout:       90 * time.Second, // 空闲连接保持时间
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}

		// 创建客户端并配置超时
		httpClient = &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second, // 总请求超时时间
		}
	})

	return httpClient
}

// doHttp 函数用于执行 HTTP 请求并解析 JSON 响应
// 参数 method: HTTP 请求方法（如 GET、POST）
// 参数 url: 请求的 URL 路径
// 参数 body: 请求体内容
// 返回值: 解析后的 JSON 对象和可能的错误
func (a *app) doHttp(method string, url string, body io.Reader) (*json2.Json, error) {
	req, err := http.NewRequest(method, a.server+url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	response, err := createHTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err)
		}
	}(response.Body)

	resp, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	js, err := json2.NewJson(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if errCode := js.Get("errcode"); errCode != nil && errCode.MustInt() != 0 {
		return js, errors.New(js.Get("errmsg").MustString())
	}

	return js, nil
}

func (a *app) Id() string {
	return a.config.AppId
}

func (a *app) Token() string {
	return a.token.GetAccessToken()
}
