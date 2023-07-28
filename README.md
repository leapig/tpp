# TPP(Third-Party Platform)

## 前言
> 支付钉钉、飞书、微信公众号、微信小程序、企业微信、微信开放平台等接入

## 快速开始 & demo

```shell script
go get github.com/leapig/tpp
```

```go
# 实例对应平台
// 企微【企业内部开发】/【服务商代开发】应用
app :=NewTpp().WW(ww.Config{
	CorpId:"企业ID"
	CorpSecret:"应用的凭证密钥"
})
app :=NewTpp().DT(dt.Config)    // 钉钉
app :=NewTpp().FS(fs.Config)    // 飞书
app :=NewTpp().MP(mp.Config)    // 微信小程序
app :=NewTpp().OA(oa.Config)    // 微信公众号
app :=NewTpp().WO(wo.Config)    // 微信开放平台
# 调用平台接口
app.DoAnything()
```