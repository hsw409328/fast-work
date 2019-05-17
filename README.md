# FastWork

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

### 整体使用

下一步版本

