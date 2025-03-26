package util

import (
	"encoding/json"
	"github.com/faabiosr/cachego"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const ContentType = "application/json;charset=utf-8"

type getRefreshRequestFunc func() []byte

type AccessToken struct {
	Id                    string
	Cache                 cachego.Cache
	GetRefreshRequestFunc getRefreshRequestFunc
}

var refreshAccessTokenLock sync.Mutex

func (a AccessToken) GetAccessToken() (token string) {
	refreshAccessTokenLock.Lock()
	defer refreshAccessTokenLock.Unlock()
	token, _ = a.Cache.Fetch("access_token:" + a.Id)
	if token != "" {
		return
	}

	resp := a.GetRefreshRequestFunc()
	var res struct {
		AccessToken           string `json:"access_token"`
		ComponentAccessToken  string `json:"component_access_token"`
		AuthorizerAccessToken string `json:"authorizer_access_token"`
		AccessDingToken       string `json:"accessToken"`
		AccessLarkToken       string `json:"tenant_access_token"`
		ExpiresIn             int    `json:"expires_in"`
		Expire                int    `json:"expire"`
	}
	_ = json.Unmarshal(resp, &res)
	if res.AccessToken != "" {
		token = res.AccessToken
		d := time.Duration(res.ExpiresIn) * time.Second
		_ = a.Cache.Save("access_token:"+a.Id, token, d)
	} else if res.ComponentAccessToken != "" {
		token = res.AuthorizerAccessToken
		d := time.Duration(res.ExpiresIn) * time.Second
		_ = a.Cache.Save("access_token:"+a.Id, token, d)
	} else if res.AuthorizerAccessToken != "" {
		token = res.AuthorizerAccessToken
		d := time.Duration(res.ExpiresIn) * time.Second
		_ = a.Cache.Save("access_token:"+a.Id, token, d)
	} else if res.AccessLarkToken != "" {
		token = res.AccessLarkToken
		d := time.Duration(res.Expire) * time.Second
		_ = a.Cache.Save("access_token:"+a.Id, token, d)
	} else if res.AccessDingToken != "" {
		token = res.AccessDingToken
		d := time.Duration(res.ExpiresIn) * time.Second
		_ = a.Cache.Save("access_token:"+a.Id, token, d)
	}
	return
}

func (a AccessToken) ApplyAccessToken(url url.Values) url.Values {
	url.Add("access_token", a.GetAccessToken())
	return url
}

func (a AccessToken) SetLarkAccessToken(header http.Header) http.Header {
	header.Set("Authorization", "Bearer "+a.GetAccessToken())
	return header
}
