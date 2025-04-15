package wo

import (
	"bytes"
	"encoding/json"
	json2 "github.com/bitly/go-simplejson"
	"net/http"
	"net/url"
)

// ThirdpartyCode2Session 小程序登录
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/login/thirdpartyCode2Session.html
// req GET https://api.weixin.qq.com/sns/component/jscode2session?component_access_token=ACCESS_TOKEN
func (a *app) ThirdpartyCode2Session(appid, jsCode string) (*json2.Json, error) {
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
func (a *app) SetPrivacySetting(authorizerAccessToken string, privacyVer int64, settingList []interface{}, ownerSettingList interface{}, sdkPrivacyInfoList []interface{}) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	payload, _ := json.Marshal(map[string]interface{}{
		"privacy_ver":           privacyVer,
		"setting_list":          settingList,
		"owner_setting":         ownerSettingList,
		"sdk_privacy_info_list": sdkPrivacyInfoList,
	})
	return a.doHttp(http.MethodPost, "/cgi-bin/component/setprivacysetting?"+params.Encode(), bytes.NewReader(payload))
}

// GetPrivacySetting 获取小程序用户隐私保护指引
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/privacy-management/getPrivacySetting.html
// req POST https://api.weixin.qq.com/cgi-bin/component/getprivacysetting?access_token=ACCESS_TOKEN
func (a *app) GetPrivacySetting(authorizerAccessToken string, privacyVer int64) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	payload, _ := json.Marshal(map[string]interface{}{
		"privacy_ver": privacyVer,
	})
	return a.doHttp(http.MethodPost, "/cgi-bin/component/getprivacysetting?"+params.Encode(), bytes.NewReader(payload))
}

// UploadPrivacySetting 上传小程序用户隐私保护指引
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/privacy-management/uploadPrivacySetting.html
// req POST https://api.weixin.qq.com/cgi-bin/component/uploadprivacyextfile?access_token=ACCESS_TOKEN
func (a *app) UploadPrivacySetting(authorizerAccessToken string, file *bytes.Buffer) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	return a.doHttp(http.MethodPost, "/cgi-bin/component/uploadprivacyextfile?"+params.Encode(), file)
}
