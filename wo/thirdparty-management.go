package wo

import (
	"bytes"
	"encoding/json"
	json2 "github.com/bitly/go-simplejson"
	"net/http"
	"net/url"
	"strings"
)

/* template-management */

// GetTemplatedRaftList 获取草稿箱列表
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/thirdparty-management/template-management/getTemplatedRaftList.html
// req GET https://api.weixin.qq.com/wxa/gettemplatedraftlist?access_token=ACCESS_TOKEN
func (a *app) GetTemplatedRaftList() (js *json2.Json, err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	req, err := http.NewRequest(http.MethodGet, a.server+"/wxa/gettemplatedraftlist?"+params.Encode(), nil)
	if err == nil {
		res, err := a.do(req)
		if err == nil {
			js = res.Get("draft_list")
		}
	}
	return
}

// AddToTemplate 将草稿添加到模板库
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/thirdparty-management/template-management/addToTemplate.html
// req POST https://api.weixin.qq.com/wxa/addtotemplate?access_token=ACCESS_TOKEN
func (a *app) AddToTemplate(draftId, templateType int64) (js *json2.Json, err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]int64{
		"draft_id":      draftId,
		"template_type": templateType,
	})
	req, err := http.NewRequest(http.MethodPost, a.server+"/wxa/addtotemplate?"+params.Encode(), bytes.NewReader(payload))
	if err == nil {
		js, err = a.do(req)
	}
	return
}

// GetTemplateList 获取模板列表
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/thirdparty-management/template-management/getTemplateList.html
// req GET https://api.weixin.qq.com/wxa/gettemplatelist?access_token=ACCESS_TOKEN
func (a *app) GetTemplateList(templateType int64) (js *json2.Json, err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]int64{
		"template_type": templateType,
	})
	req, err := http.NewRequest(http.MethodGet, a.server+"/wxa/gettemplatelist?"+params.Encode(), bytes.NewReader(payload))
	if err == nil {
		js, err = a.do(req)
	}
	return
}

// DeleteTemplate 删除代码模板
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/thirdparty-management/template-management/deleteTemplate.html
// req POST https://api.weixin.qq.com/wxa/deletetemplate?access_token=ACCESS_TOKEN
func (a *app) DeleteTemplate(templateId int64) (js *json2.Json, err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	payload, _ := json.Marshal(map[string]int64{
		"template_id": templateId,
	})
	req, err := http.NewRequest(http.MethodPost, a.server+"/wxa/deletetemplate?"+params.Encode(), bytes.NewReader(payload))
	if err == nil {
		js, err = a.do(req)
	}
	return
}

/* domain-mgnt */

// ModifyThirdpartyServerDomain 设置第三方平台服务器域名
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/thirdparty-management/domain-mgnt/modifyThirdpartyServerDomain.html
// req POST https://api.weixin.qq.com/cgi-bin/component/modify_wxa_server_domain?access_token=ACCESS_TOKEN
func (a *app) ModifyThirdpartyServerDomain(action, WxaServerDomain string, IsModifyPublishedTogether bool) (js *json2.Json, err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	var body map[string]interface{}
	if strings.ToLower(action) != "get" {
		body = map[string]interface{}{
			"action":                       action,
			"wxa_server_domain":            WxaServerDomain,
			"is_modify_published_together": IsModifyPublishedTogether,
		}
	} else {
		body = map[string]interface{}{
			"action": action,
		}
	}
	payload, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/modify_wxa_server_domain?"+params.Encode(), bytes.NewReader(payload))
	if err == nil {
		js, err = a.do(req)
	}
	return
}

// GetThirdpartyJumpDomainConfirmFile 获取第三方平台业务域名校验文件
// doc https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/thirdparty-management/domain-mgnt/getThirdpartyJumpDomainConfirmFile.html
// req POST https://api.weixin.qq.com/cgi-bin/component/get_domain_confirmfile?access_token=ACCESS_TOKEN
func (a *app) GetThirdpartyJumpDomainConfirmFile() (js *json2.Json, err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	req, err := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/get_domain_confirmfile?"+params.Encode(), nil)
	if err == nil {
		js, err = a.do(req)
	}
	return
}

// ModifyThirdpartyJumpDomain 设置第三方平台业务域名
// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/thirdparty-management/domain-mgnt/modifyThirdpartyJumpDomain.html
// req POST https://api.weixin.qq.com/cgi-bin/component/modify_wxa_jump_domain?access_token=ACCESS_TOKEN
func (a *app) ModifyThirdpartyJumpDomain(action, WxaJumpH5Domain string, IsModifyPublishedTogether bool) (js *json2.Json, err error) {
	params := url.Values{}
	params = a.token.ApplyAccessToken(params)
	var body map[string]interface{}
	if strings.ToLower(action) != "get" {
		body = map[string]interface{}{
			"action":                       action,
			"wxa_jump_h5_domain":           WxaJumpH5Domain,
			"is_modify_published_together": IsModifyPublishedTogether,
		}
	} else {
		body = map[string]interface{}{
			"action": action,
		}
	}
	payload, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, a.server+"/cgi-bin/component/modify_wxa_jump_domain?"+params.Encode(), bytes.NewReader(payload))
	if err == nil {
		js, err = a.do(req)
	}
	return
}
