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
func (a *app) SetPrivacySetting(authorizerAccessToken string, privacyVer int64, settingList, ownerSettingList, sdkPrivacyInfoList interface{}) (*json2.Json, error) {
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

// TODO

// SubmitAudit 提交代码审核
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/submitAudit.html
// req POST https://api.weixin.qq.com/wxa/submit_audit?access_token=ACCESS_TOKEN
func (a *app) SubmitAudit(authorizerAccessToken string, itemList interface{}, feedbackInfo, feedbackStuff, versionDesc string, previewInfo map[string]interface{}, ugcDeclare map[string]interface{}, privacyApiNotUse bool, orderPath string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	payload, _ := json.Marshal(map[string]interface{}{
		"item_list":           itemList,
		"feedback_info":       feedbackInfo,
		"feedback_stuff":      feedbackStuff,
		"version_desc":        versionDesc,
		"preview_info":        previewInfo,
		"ugc_declare":         ugcDeclare,
		"privacy_api_not_use": privacyApiNotUse,
		"order_path":          orderPath,
	})
	return a.doHttp(http.MethodPost, "/wxa/submit_audit?"+params.Encode(), bytes.NewReader(payload))
}

// GetAuditStatus 查询审核单状态
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/getAuditStatus.html
// req POST https://api.weixin.qq.com/wxa/get_auditstatus?access_token=ACCESS_TOKEN
func (a *app) GetAuditStatus(authorizerAccessToken string, auditId int64) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	payload, _ := json.Marshal(map[string]interface{}{
		"auditid": auditId,
	})
	return a.doHttp(http.MethodPost, "/wxa/get_auditstatus?"+params.Encode(), bytes.NewReader(payload))
}

// UndoAudit 撤回代码审核
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/undoAudit.html
// req GET https://api.weixin.qq.com/wxa/undocodeaudit?access_token=ACCESS_TOKEN
func (a *app) UndoAudit(authorizerAccessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	return a.doHttp(http.MethodGet, "/wxa/undocodeaudit?"+params.Encode(), nil)
}

// Release 发布已通过审核的小程序
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/release.html
// req POST https://api.weixin.qq.com/wxa/release?access_token=ACCESS_TOKEN
func (a *app) Release(authorizerAccessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	payload, _ := json.Marshal(map[string]interface{}{})
	return a.doHttp(http.MethodPost, "/wxa/release?"+params.Encode(), bytes.NewReader(payload))
}

// GetVersionInfo 查询小程序版本信息
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/getVersionInfo.html
// POST https://api.weixin.qq.com/wxa/getversioninfo?access_token=ACCESS_TOKEN
func (a *app) GetVersionInfo(authorizerAccessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	payload, _ := json.Marshal(map[string]interface{}{})
	return a.doHttp(http.MethodPost, "/wxa/getversioninfo?"+params.Encode(), bytes.NewBuffer(payload))
}

// GetLatestAuditStatus 查询最新一次提交的审核状态
// doc https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/code/get_latest_auditstatus.html
// req GET https://api.weixin.qq.com/wxa/get_latest_auditstatus?access_token=ACCESS_TOKEN
func (a *app) GetLatestAuditStatus(authorizerAccessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	return a.doHttp(http.MethodGet, "/wxa/get_latest_auditstatus?"+params.Encode(), nil)
}

// UploadMediaToCodeAudit 上传提审素材
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/uploadMediaToCodeAudit.html
// req POST https://api.weixin.qq.com/wxa/uploadmedia?access_token=ACCESS_TOKEN
func (a *app) UploadMediaToCodeAudit(authorizerAccessToken string, file *bytes.Buffer) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	return a.doHttp(http.MethodPost, "/wxa/uploadmedia?"+params.Encode(), file)
}

// GetCodePrivacyInfo 获取隐私接口检测结果
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/getCodePrivacyInfo.html
// req GET https://api.weixin.qq.com/wxa/security/get_code_privacy_info?access_token=ACCESS_TOKEN
func (a *app) GetCodePrivacyInfo(authorizerAccessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	return a.doHttp(http.MethodGet, "/wxa/security/get_code_privacy_info?"+params.Encode(), nil)
}
