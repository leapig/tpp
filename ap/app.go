package ap

import (
	"context"
	"github.com/go-pay/gopay/alipay"
)

type App interface {
	SystemOauthToken(code string) map[string]interface{}
}

type Config struct {
	AppId      string `json:"appid"`
	AesKey     string `json:"aesKey"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type app struct {
	config Config
	server string
}

func NewApp(config Config) App {
	return &app{
		config: config,
	}
}

// SystemOauthToken 获取用户登录信息
func (a *app) SystemOauthToken(code string) (res map[string]interface{}) {
	if resp, err := alipay.SystemOauthToken(
		context.Background(),
		a.config.AppId,
		a.config.PrivateKey,
		"authorization_code",
		code, "RSA2"); err == nil && resp.ErrorResponse == nil {
		res = map[string]interface{}{
			"openid":  resp.Response.UserId,
			"unionid": resp.Response.UnionId,
		}
	}
	return
}
