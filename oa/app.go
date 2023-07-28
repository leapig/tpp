package oa

type App interface {
	Config() App
}

type Config struct {
}

type app struct {
	config Config
}

func NewApp(config Config) App {
	return &app{config: config}
}

func (a *app) Config() App {
	return a
}
