package wo

import (
	json2 "github.com/bitly/go-simplejson"
	"net/http"
	"net/url"
)

// ThirdpartyCode2Session 小程序登录
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/login/thirdpartyCode2Session.html
// req GET https://api.weixin.qq.com/sns/component/jscode2session?component_access_token=ACCESS_TOKEN
func (a *app) ThirdpartyCode2Session(appid, jsCode string) (js *json2.Json, err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	params.Add("appid", appid)
	params.Add("js_code", jsCode)
	params.Add("grant_type", "authorization_code")
	params.Add("component_appid", a.config.AppId)
	return a.doHttp(http.MethodGet, "/sns/component/jscode2session?"+params.Encode(), nil)
}

// TODO

// SetPrivacySetting 设置小程序用户隐私保护指引
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/privacy-management/setPrivacySetting.html
// req POST https://api.weixin.qq.com/cgi-bin/component/setprivacysetting?access_token=ACCESS_TOKEN
