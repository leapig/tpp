package tpp

import (
	"tpp/dt"
	"tpp/fs"
	"tpp/mp"
	"tpp/oa"
	"tpp/wo"
	"tpp/ww"
)

type Tpp interface {
	WW(ww.Config) ww.App
	DT(dt.Config) dt.App
	FS(fs.Config) fs.App
	MP(mp.Config) mp.App
	OA(oa.Config) oa.App
	WO(wo.Config) wo.App
}

type tpp struct{}

func NewTpp() Tpp {
	return &tpp{}
}

func (t *tpp) WW(config ww.Config) ww.App {
	return ww.NewApp(config)
}
func (t *tpp) DT(config dt.Config) dt.App {
	return dt.NewApp(config)
}
func (t *tpp) FS(config fs.Config) fs.App {
	return fs.NewApp(config)
}
func (t *tpp) MP(config mp.Config) mp.App {
	return mp.NewApp(config)
}
func (t *tpp) OA(config oa.Config) oa.App {
	return oa.NewApp(config)
}
func (t *tpp) WO(config wo.Config) wo.App {
	return wo.NewApp(config)
}
