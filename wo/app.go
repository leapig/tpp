package wo

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
	"os"
	"sync"
	"time"
)

type App interface {
	// Id 获取AppId
	Id() string
	// Token 获取Token
	Token() string
	// GetAuthorizerList 拉取已授权的账号信息
	// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/authorization-management/getAuthorizerList.html
	// req POST https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_list?access_token=ACCESS_TOKEN
	GetAuthorizerList() ([]*json2.Json, error)
	// GetAuthorizerInfo 获取授权账号详情
	// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/authorization-management/getAuthorizerInfo.html
	// req POST https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info?access_token=ACCESS_TOKEN
	GetAuthorizerInfo(authorizerAppId string) (*json2.Json, error)
	// SetAuthorizerOptionInfo 设置授权方选项信息
	// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/authorization-management/setAuthorizerOptionInfo.html
	// req POST https://api.weixin.qq.com/cgi-bin/component/set_authorizer_option?access_token=ACCESS_TOKEN
	SetAuthorizerOptionInfo(authorizerAccessToken, optionName, optionValue string) (*json2.Json, error)
	// GetAuthorizerOptionInfo 获取授权方选项信息
	// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/authorization-management/getAuthorizerOptionInfo.html
	// req POST https://api.weixin.qq.com/cgi-bin/component/get_authorizer_option?access_token=ACCESS_TOKEN
	GetAuthorizerOptionInfo(authorizerAccessToken, optionName string) (*json2.Json, error)
	// GetTemplatedRaftList 获取草稿箱列表
	// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/thirdparty-management/template-management/getTemplatedRaftList.html
	// req GET https://api.weixin.qq.com/wxa/gettemplatedraftlist?access_token=ACCESS_TOKEN
	GetTemplatedRaftList() (*json2.Json, error)
	// AddToTemplate 将草稿添加到模板库
	// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/thirdparty-management/template-management/addToTemplate.html
	// req POST https://api.weixin.qq.com/wxa/addtotemplate?access_token=ACCESS_TOKEN
	AddToTemplate(draftId, templateType int64) (*json2.Json, error)
	// GetTemplateList 获取模板列表
	// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/thirdparty-management/template-management/getTemplateList.html
	// req GET https://api.weixin.qq.com/wxa/gettemplatelist?access_token=ACCESS_TOKEN
	GetTemplateList(templateType int64) (*json2.Json, error)
	// DeleteTemplate 删除代码模板
	// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/thirdparty-management/template-management/deleteTemplate.html
	// req POST https://api.weixin.qq.com/wxa/deletetemplate?access_token=ACCESS_TOKEN
	DeleteTemplate(templateId int64) (*json2.Json, error)
	// ModifyThirdpartyServerDomain 设置第三方平台服务器域名
	// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/thirdparty-management/domain-mgnt/modifyThirdpartyServerDomain.html
	// req POST https://api.weixin.qq.com/cgi-bin/component/modify_wxa_server_domain?access_token=ACCESS_TOKEN
	ModifyThirdpartyServerDomain(action, WxaServerDomain string, IsModifyPublishedTogether bool) (*json2.Json, error)
	// GetThirdpartyJumpDomainConfirmFile 获取第三方平台业务域名校验文件
	// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/thirdparty-management/domain-mgnt/getThirdpartyJumpDomainConfirmFile.html
	// req POST https://api.weixin.qq.com/cgi-bin/component/get_domain_confirmfile?access_token=ACCESS_TOKEN
	GetThirdpartyJumpDomainConfirmFile() (js *json2.Json, err error)
	// ModifyThirdpartyJumpDomain 设置第三方平台业务域名
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/thirdparty-management/domain-mgnt/modifyThirdpartyJumpDomain.html
	// req POST https://api.weixin.qq.com/cgi-bin/component/modify_wxa_jump_domain?access_token=ACCESS_TOKEN
	ModifyThirdpartyJumpDomain(action, WxaJumpH5Domain string, IsModifyPublishedTogether bool) (*json2.Json, error)
	// StartPushTicket 开启推送ticket
	// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/startPushTicket.html
	// req POST https://api.weixin.qq.com/cgi-bin/component/api_start_push_ticket
	StartPushTicket() (*json2.Json, error)
	// GetPreAuthCode 获取预授权码
	// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getPreAuthCode.html
	// req POST https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?access_token=ACCESS_TOKEN
	GetPreAuthCode() (*json2.Json, error)
	// GetAuthorizerAccessToken 获取授权账号调用令牌
	// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getAuthorizerAccessToken.html
	// req POST https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token=ACCESS_TOKEN
	GetAuthorizerAccessToken(authorizerAppId, authorizerRefreshToken string) (*json2.Json, error)
	// GetAuthorizerRefreshToken 获取刷新令牌
	// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getAuthorizerRefreshToken.html
	// req POST https://api.weixin.qq.com/cgi-bin/component/api_query_auth?access_token=ACCESS_TOKEN
	GetAuthorizerRefreshToken(authorizationCode string) (*json2.Json, error)
	// GetComponentAccessToken 获取令牌
	// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getComponentAccessToken.html
	// req POST https://api.weixin.qq.com/cgi-bin/component/api_component_token
	GetComponentAccessToken() (*json2.Json, error)
}

type Config struct {
	AppId  string        `json:"appid"`
	Secret string        `json:"secret"`
	Token  string        `json:"token"`
	AesKey string        `json:"aes_key"`
	Ticket string        `json:"ticket"`
	Cache  cachego.Cache `json:"cache"`
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
