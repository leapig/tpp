package wo

import (
	"bytes"
	"encoding/json"
	"errors"
	json2 "github.com/bitly/go-simplejson"
	"io"
	"net/http"
	"net/url"
	"strings"
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

// GetAccountBasicInfo 获取基本信息
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/basic-info-management/getAccountBasicInfo.html
// req POST https://api.weixin.qq.com/cgi-bin/account/getaccountbasicinfo?access_token=ACCESS_TOKEN
func (a *app) GetAccountBasicInfo(authorizerAccessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	return a.doHttp(http.MethodPost, "/cgi-bin/account/getaccountbasicinfo?"+params.Encode(), nil)
}

// GetBindOpenAccount 查询绑定的开放平台账号
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/basic-info-management/getBindOpenAccount.html
// req GET https://api.weixin.qq.com/cgi-bin/open/have?access_token=ACCESS_TOKEN
func (a *app) GetBindOpenAccount(authorizerAccessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	return a.doHttp(http.MethodGet, "/cgi-bin/open/have?"+params.Encode(), nil)
}

// ModifyServerDomain 配置小程序服务器域名
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/domain-management/modifyServerDomain.html
// req POST https://api.weixin.qq.com/wxa/modify_domain?access_token=ACCESS_TOKEN
func (a *app) ModifyServerDomain(authorizerAccessToken, action string, requestDomain, wsRequestDomain, uploadDomain, downloadDomain, udpDomain, tcpDomain []string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	var body map[string]interface{}
	if strings.ToLower(action) != "get" {
		body = map[string]interface{}{
			"action":          action,
			"requestdomain":   requestDomain,
			"wsrequestdomain": wsRequestDomain,
			"uploaddomain":    uploadDomain,
			"downloaddomain":  downloadDomain,
			"udpdomain":       udpDomain,
			"tcpdomain":       tcpDomain,
		}
	} else {
		body = map[string]interface{}{
			"action": action,
		}
	}
	payload, _ := json.Marshal(body)
	return a.doHttp(http.MethodPost, "/wxa/modify_domain?"+params.Encode(), bytes.NewReader(payload))
}

// ModifyJumpDomain 配置小程序业务域名
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/domain-management/modifyJumpDomain.html
// req POST https://api.weixin.qq.com/wxa/setwebviewdomain?access_token=ACCESS_TOKEN
func (a *app) ModifyJumpDomain(authorizerAccessToken, action string, webviewDomain []string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	var body map[string]interface{}
	if strings.ToLower(action) != "get" {
		body = map[string]interface{}{
			"action":        action,
			"webviewdomain": webviewDomain,
		}
	} else {
		body = map[string]interface{}{
			"action": action,
		}
	}
	payload, _ := json.Marshal(body)
	return a.doHttp(http.MethodPost, "/wxa/setwebviewdomain?"+params.Encode(), bytes.NewReader(payload))
}

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

// Commit 上传代码并生成体验版
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/commit.html
// req POST https://api.weixin.qq.com/wxa/commit?access_token=ACCESS_TOKEN
func (a *app) Commit(authorizerAccessToken, templateId, extJson, userVersion, userDesc string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	payload, _ := json.Marshal(map[string]interface{}{
		"template_id":  templateId,
		"ext_json":     extJson,
		"user_version": userVersion,
		"user_desc":    userDesc,
	})
	return a.doHttp(http.MethodPost, "/wxa/commit?"+params.Encode(), bytes.NewReader(payload))
}

// GetCodePage 获取已上传的代码页面列表
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/getCodePage.html
// req GET https://api.weixin.qq.com/wxa/get_page?access_token=ACCESS_TOKEN
func (a *app) GetCodePage(authorizerAccessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	return a.doHttp(http.MethodGet, "/wxa/commit?"+params.Encode(), nil)
}

// GetTrialQRCode 获取体验版二维码
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/getTrialQRCode.html
// req GET https://api.weixin.qq.com/wxa/get_qrcode?access_token=ACCESS_TOKEN
func (a *app) GetTrialQRCode(authorizerAccessToken, path string) ([]byte, error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	params.Add("path", url.QueryEscape(path))
	if response, err := http.Get(a.server + "/wxa/get_qrcode?" + params.Encode()); err == nil {
		if resp, err := io.ReadAll(response.Body); err == nil {
			if response.Header.Get("Content-Type") == "image/jpeg" {
				return resp, err
			} else {
				js, err := json2.NewJson(resp)
				if err != nil {
					return nil, err
				}
				if js.Get("errcode").MustInt() != 0 {
					return nil, errors.New(js.Get("errmsg").MustString())
				}
				return nil, nil
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

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

// RevertCodeReleaseGetVersion 小程序版本回退(获取可回退的小程序版本)
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/revertCodeRelease.html
// req GET https://api.weixin.qq.com/wxa/revertcoderelease?access_token=ACCESS_TOKEN
func (a *app) RevertCodeReleaseGetVersion(authorizerAccessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	params.Add("action", "get_history_version")
	return a.doHttp(http.MethodGet, "/wxa/revertcoderelease?"+params.Encode(), nil)
}

// RevertCodeReleaseRollback 小程序版本回退(回滚到指定的小程序版本，默认上一个版本)
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/revertCodeRelease.html
// req GET https://api.weixin.qq.com/wxa/revertcoderelease?access_token=ACCESS_TOKEN
func (a *app) RevertCodeReleaseRollback(authorizerAccessToken, appVersion string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	if appVersion != "" {
		params.Add("app_version", appVersion)
	}
	return a.doHttp(http.MethodGet, "/wxa/revertcoderelease?"+params.Encode(), nil)
}

// GrayRelease 分阶段发布
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/grayRelease.html
// req POST https://api.weixin.qq.com/wxa/grayrelease?access_token=ACCESS_TOKEN
func (a *app) GrayRelease(authorizerAccessToken string, grayPercentage int64, supportDebugerFirst, supportExperiencerFirst bool) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	payload, _ := json.Marshal(map[string]interface{}{
		"gray_percentage":           grayPercentage,
		"support_debuger_first":     supportDebugerFirst,
		"support_experiencer_first": supportExperiencerFirst,
	})
	return a.doHttp(http.MethodPost, "/wxa/grayrelease?"+params.Encode(), bytes.NewReader(payload))
}

// GetGrayReleasePlan 获取分阶段发布详情
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/getGrayReleasePlan.html
// req GET https://api.weixin.qq.com/wxa/getgrayreleaseplan?access_token=ACCESS_TOKEN
func (a *app) GetGrayReleasePlan(authorizerAccessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	return a.doHttp(http.MethodGet, "/wxa/getgrayreleaseplan?"+params.Encode(), nil)
}

// SetVisitStatus 设置小程序服务状态
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/setVisitStatus.html
// req POST https://api.weixin.qq.com/wxa/change_visitstatus?access_token=ACCESS_TOKEN
func (a *app) SetVisitStatus(authorizerAccessToken string, action string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	payload, _ := json.Marshal(map[string]interface{}{
		"action": action,
	})
	return a.doHttp(http.MethodPost, "/wxa/change_visitstatus?"+params.Encode(), bytes.NewReader(payload))
}

// RevertGrayRelease 取消分阶段发布
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/revertGrayRelease.html
// req GET https://api.weixin.qq.com/wxa/revertgrayrelease?access_token=ACCESS_TOKEN
func (a *app) RevertGrayRelease(authorizerAccessToken string) (*json2.Json, error) {
	params := url.Values{}
	params.Add("access_token", authorizerAccessToken)
	return a.doHttp(http.MethodGet, "/wxa/revertgrayrelease?"+params.Encode(), nil)
}

// GetVersionInfo 查询小程序版本信息
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/getVersionInfo.html
// req POST https://api.weixin.qq.com/wxa/getversioninfo?access_token=ACCESS_TOKEN
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
