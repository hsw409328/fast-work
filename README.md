# FastWork
快速收集资产，快速应用于企业安全项目。简单化操作，支持子项目模块独立使用，可以自由进行拆分。
抓取网页使用的是chrome浏览器内核，支持单页应用，区别于普通的爬虫项目。

### 单个项目使用
##### 爆破子域名
```
object := fast_dns_search.NewDnsBlast("xxxx.com")
object.A = true
object.Run()

#结果
存入redis 请先配置本地redis
```

##### 抓取网页
```
fast_crawl_engine.NewFastCrawlEngine(FastCrawlEngineParams{
    BaseDomain:   "http://www.51hsw.com",
    DomainStr:    "http://www.51hsw.com",
    MinDeepLevel: 1,
    MaxDeepLevel: 2,
}).Start()

#结果

```
### 新增
* 新增过滤器
* 爬虫服务端

### 下一步版本
* WEB管理平台
* 爬虫分布式部署
* ...

### 相关项目

* chrome浏览器内核
* https://github.com/chromedp/chromedp
* https://github.com/benmanns/goworker
* https://github.com/hsw409328/gofunc
*



