package tpp

import (
	"tpp/ap"
	"tpp/dt"
	"tpp/fs"
	"tpp/mp"
	"tpp/oa"
	"tpp/wk"
	"tpp/wo"
	"tpp/ww"
)

type Tpp interface {
	WW(ww.Config) ww.App
	MP(mp.Config) mp.App
	OA(oa.Config) oa.App
	WK(wk.Config) wk.App
	WO(wo.Config) wo.App
	AP(ap.Config) ap.App
	DT(dt.Config) dt.App
	FS(fs.Config) fs.App
}

type tpp struct{}

func NewTpp() Tpp {
	return &tpp{}
}

// WW 企业微信实例
func (t *tpp) WW(config ww.Config) ww.App {
	return ww.NewApp(config)
}

// MP 微信小程序实例
func (t *tpp) MP(config mp.Config) mp.App {
	return mp.NewApp(config)
}

// OA 微信公众号实例
func (t *tpp) OA(config oa.Config) oa.App {
	return oa.NewApp(config)
}

// WK 腾讯微卡实例
func (t *tpp) WK(config wk.Config) wk.App {
	return wk.NewApp(config)
}

// WO 微信开放平台实例
func (t *tpp) WO(config wo.Config) wo.App {
	return wo.NewApp(config)
}

// AP 支付宝实例
func (t *tpp) AP(config ap.Config) ap.App {
	return ap.NewApp(config)
}

// DT 钉钉实例
func (t *tpp) DT(config dt.Config) dt.App {
	return dt.NewApp(config)
}

// FS 飞书实例
func (t *tpp) FS(config fs.Config) fs.App {
	return fs.NewApp(config)
}
