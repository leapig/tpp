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
func (a *app) ClearQuota(appId, accessToken string) (js *json2.Json, err error) {
	params := url.Values{}
	params.Add("access_token", accessToken)
	payload, _ := json.Marshal(map[string]string{
		"appid": appId,
	})
	req, err := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/clear_quota?"+params.Encode(), bytes.NewReader(payload))
	if err == nil {
		js, err = a.do(req)
	}
	return
}

// GetApiQuota 查询API调用额度
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/openapi/getApiQuota.html
// req POST https://api.weixin.qq.com/cgi-bin/openapi/quota/get?access_token=ACCESS_TOKEN
func (a *app) GetApiQuota(cgiPath, accessToken string) (js *json2.Json, err error) {
	params := url.Values{}
	params.Add("access_token", accessToken)
	payload, _ := json.Marshal(map[string]string{
		"cgi_path": cgiPath,
	})
	req, err := http.NewRequest(http.MethodPost, a.server+"/openapi/quota/get?"+params.Encode(), bytes.NewReader(payload))
	if err == nil {
		js, err = a.do(req)
	}
	return
}

// GetRidInfo 查询rid信息
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/openapi/getRidInfo.html
// req POST https://api.weixin.qq.com/cgi-bin/openapi/rid/get?access_token=ACCESS_TOKEN

// ClearComponentQuotaByAppSecret 使用AppSecret重置第三方平台 API 调用次数
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/openapi/clearComponentQuotaByAppSecret.html
// req POST https://api.weixin.qq.com/cgi-bin/component/clear_quota/v2
