package wo

import (
	"bytes"
	"encoding/json"
	json2 "github.com/bitly/go-simplejson"
	"net/http"
	"net/url"
)

// ClearQuota 重置API调用次数
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/openapi/clearQuota.html
// req POST https://api.weixin.qq.com/cgi-bin/clear_quota?access_token=ACCESS_TOKEN
func (a *app) ClearQuota(appId, accessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", accessToken)
	payload, _ := json.Marshal(map[string]string{
		"appid": appId,
	})
	return a.doHttp(http.MethodPost, "/cgi-bin/clear_quota?"+params.Encode(), bytes.NewReader(payload))
}

// GetApiQuota 查询API调用额度
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/openapi/getApiQuota.html
// req POST https://api.weixin.qq.com/cgi-bin/openapi/quota/get?access_token=ACCESS_TOKEN
func (a *app) GetApiQuota(cgiPath, accessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", accessToken)
	payload, _ := json.Marshal(map[string]string{
		"cgi_path": cgiPath,
	})
	return a.doHttp(http.MethodPost, "/cgi-bin/openapi/quota/get?"+params.Encode(), bytes.NewReader(payload))
}

// GetRidInfo 查询rid信息
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/openapi/getRidInfo.html
// req POST https://api.weixin.qq.com/cgi-bin/openapi/rid/get?access_token=ACCESS_TOKEN
func (a *app) GetRidInfo(rid, accessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", accessToken)
	payload, _ := json.Marshal(map[string]string{
		"rid": rid,
	})
	return a.doHttp(http.MethodPost, "/cgi-bin/openapi/rid/get?"+params.Encode(), bytes.NewReader(payload))
}

// ClearComponentQuotaByAppSecret 使用AppSecret重置第三方平台 API 调用次数
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/openapi/clearComponentQuotaByAppSecret.html
// req POST https://api.weixin.qq.com/cgi-bin/component/clear_quota/v2
func (a *app) ClearComponentQuotaByAppSecret(appid string) (*json2.Json, error) {
	body := map[string]string{
		"appid":           appid,
		"component_appid": a.config.AppId,
		"appsecret":       a.config.Secret,
	}
	if appid == a.config.AppId {
		delete(body, "appid")
	}
	payload, _ := json.Marshal(body)
	return a.doHttp(http.MethodPost, "/cgi-bin/component/clear_quota/v2", bytes.NewReader(payload))
}
