/**
 * Author: haoshuaiwei 
 * Date: 2019-05-14 18:14 
 */

package fast_driver

import (
	"fast-work/fast-sys"
	"gopkg.in/redis.v5"
	"strconv"
	"time"
)

/**
存储已经爆破的域名结构使用
dns-search-list
	[jj.com,aa.com,cc.com,dd.com]

待爆破的域名队列
dns-search-wait-list
	[jj.com,aa.com,cc.com,dd.com]

需要关闭爆破的域名队列
dns-search-wait-close-list
	[jj.com,aa.com,cc.com,dd.com]


域名爆破结果结构使用 根域名+list+success
例如：
cc.com-dns-search-list-success
	{
		BaseDomain string
		Domain     string
		Ip         []string
		Host       string
		IsOpen     bool   `是否可以打开`
		Type       string `结果类型： A\TXT\CNAME\IP`
	}

爬虫结果存储结构使用 具体域名+list+success
例如：
xx.cc.com-crawl-list-success
	{
		UrlStr string
		Title string
	}
 */

const (
	RedisDnsBlastKey                = "dns-search-list"
	RedisWaitBlastKey               = "dns-search-wait-list"
	RedisWaitCloseBlastKey          = "dns-search-wait-close-list"
	RedisDomainBlastSuffixSymbolKey = "-dns-search-list-success"
)

var (
	RedisDriver *redis.Client
)

func initRedis() {
	host, _ := fast_sys.GoConfig.GetString("redis", "host")
	port, _ := fast_sys.GoConfig.GetString("redis", "port")
	authPassword, _ := fast_sys.GoConfig.GetString("redis", "auth_password")
	dbStr, _ := fast_sys.GoConfig.GetString("redis", "db")
	db, _ := strconv.Atoi(dbStr)
	RedisDriver = redis.NewClient(&redis.Options{
		Addr:        host + ":" + port,
		Password:    authPassword,
		DB:          db,
		DialTimeout: time.Second * 2,
		//IdleTimeout: time.Second * 1000000,
	})
}
