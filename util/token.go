package util

import (
	"encoding/json"
	"github.com/faabiosr/cachego"
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
		Errcode     int    `json:"errcode"`
		Errmsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	_ = json.Unmarshal(resp, &res)
	if res.Errcode == 0 {
		token = res.AccessToken
		d := time.Duration(res.ExpiresIn) * time.Second
		_ = a.Cache.Save("access_token:"+a.Id, token, d)
	}
	return
}

func (a AccessToken) ApplyAccessToken(url url.Values) url.Values {
	url.Add("access_token", a.GetAccessToken())
	return url
}
