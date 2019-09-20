/**
 * Author: haoshuaiwei 
 * Date: 2019-09-17 16:52 
 */

package controller_v1

import (
	"encoding/json"
	"fast-work/fast-crawl/fast-crawl-engine"
	"fast-work/fast-dns-search"
	"fast-work/fast-driver"
	"fast-work/fast-web/controller"
	"github.com/gin-gonic/gin"
	"github.com/hsw409328/gofunc"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type IndexController struct {
	controller.BaseController
}

// 功能导航页
func (c *IndexController) Get(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "v1-base-domain/index.html", gin.H{})
}

// 爆破根域名列表
func (c *IndexController) GetJson(ctx *gin.Context) {
	r := fast_driver.RedisDriver.LRange(fast_driver.RedisDnsBlastKey, 0, -1).Val()
	var tmp = make([]struct {
		Domain string
	}, 0)
	for _, v := range r {
		tmp = append(tmp, struct{ Domain string }{Domain: v})
	}
	ctx.JSON(http.StatusOK, c.JsonEncode(0, "success", tmp, len(tmp)))
}

// true 代表存在
func (c *IndexController) CheckDnsBlastDomainIsExist(domainStr string) bool {
	r := fast_driver.RedisDriver.LRange(fast_driver.RedisDnsBlastKey, 0, -1).Val()
	for _, v := range r {
		if domainStr == v {
			return true
		}
	}
	return false
}

// 开始搜索
func (c *IndexController) researchDns(domainStr string) {
	// 添加到待关闭爆破的队列 清除掉原来的队列
	fast_driver.RedisDriver.RPush(fast_driver.RedisWaitCloseBlastKey, domainStr)
	// 重新添加到待爆破的队列
	fast_driver.RedisDriver.RPush(fast_driver.RedisWaitBlastKey, domainStr)
	// 清空原来的爆破结果
	fast_driver.RedisDriver.Del(domainStr + fast_driver.RedisDomainBlastSuffixSymbolKey)
}

// 搜索页
func (c *IndexController) DnsSearch(ctx *gin.Context) {
	domainStr := ctx.Query("domainStr")
	isReload := ctx.Query("isReload")
	if domainStr != "" {
		err := c.CheckDnsBlastDomainIsExist(domainStr)
		if !err {
			fast_driver.RedisDriver.RPush(fast_driver.RedisDnsBlastKey, domainStr)
		}
		if isReload == "1" {
			c.researchDns(domainStr)
		}
	}
	ctx.HTML(http.StatusOK, "v1-base-domain/dns-search.html", gin.H{
		"baseDomain": domainStr,
		"isReload":   isReload,
	})
}

// 某个域名爆破结果列表
func (c *IndexController) DnsSearchJson(ctx *gin.Context) {
	domainStr := ctx.Query("domainStr")
	start, _ := strconv.Atoi(ctx.Query("start"))
	r := fast_driver.RedisDriver.LRange(domainStr+fast_driver.RedisDomainBlastSuffixSymbolKey, int64(start), -1).Val()
	var tmp = make([]fast_dns_search.DnsResultData, 0)
	var rr fast_dns_search.DnsResultData
	for _, v := range r {
		err := json.Unmarshal([]byte(v), &rr)
		if err != nil {
			log.Println(err)
			continue
		}
		tmp = append(tmp, rr)
	}
	ctx.JSON(http.StatusOK, c.JsonEncode(0, gofunc.Md5Encrypt(strings.Join(r, "")), tmp, len(tmp)))
}

// 关闭某个域名爆破结果列表
func (c *IndexController) DnsSearchClose(ctx *gin.Context) {
	domainStr := ctx.Query("domainStr")
	// 添加到待关闭爆破的队列 清除掉原来的队列
	fast_driver.RedisDriver.RPush(fast_driver.RedisWaitCloseBlastKey, domainStr)
	ctx.JSON(http.StatusOK, c.JsonEncode(0, "success", nil, 0))
}

// 爬虫
func (c *IndexController) CrawlSearch(ctx *gin.Context) {
	domainStr := ctx.Query("domainStr")
	hostStr := ctx.PostForm("host")
	maxDeepLevel, _ := strconv.Atoi(ctx.PostForm("maxDeepInt"))

	cookieValue := ctx.PostForm("cookieValue")
	cookieDomain := ctx.PostForm("cookieDomain")
	cookiePath := ctx.PostForm("cookiePath")

	isReload := ctx.Query("isReload")
	if isReload == "1" {
		if domainStr == "" || maxDeepLevel <= 0 {
			ctx.HTML(http.StatusOK, "v1-base-domain/crawl-search.html", gin.H{
				"baseDomain": domainStr,
				"isReload":   isReload,
				"err":        "添加爬虫任务失败，参数设置错误。",
			})
			return
		}
		taskObject := fast_crawl_engine.FastCrawlEngineParams{
			BaseDomain:   domainStr,
			DomainStr:    domainStr,
			MinDeepLevel: 1,
			MaxDeepLevel: maxDeepLevel,
			//Cookies: &fast_crawl_engine.FastCrawlCookies{
			//	Value:  "请使用自己的cookie",
			//	Domain: ".xxx.com",
			//	Path:   "/",
			//},
			//Cookies: nil,
			//Host: "127.0.0.1",
		}
		if cookieValue != "" && cookieDomain != "" && cookiePath != "" {
			taskObject.Cookies = &fast_crawl_engine.FastCrawlCookies{
				Value:  cookieValue,
				Domain: cookieDomain,
				Path:   cookiePath,
			}
		}
		if hostStr != "" {
			taskObject.Host = hostStr
		}
		by, err := json.Marshal(taskObject)
		if err != nil {
			ctx.HTML(http.StatusOK, "v1-base-domain/crawl-search.html", gin.H{
				"baseDomain": domainStr,
				"isReload":   isReload,
				"err":        "添加爬虫任务失败，无法开启扫描。" + err.Error(),
			})
			return
		}
		fast_driver.RedisDriver.RPush(fast_driver.RedisWaitCrawlKey, string(by))
	}
	ctx.HTML(http.StatusOK, "v1-base-domain/crawl-search.html", gin.H{
		"baseDomain": domainStr,
		"isReload":   isReload,
	})
}

// 爬虫
func (c *IndexController) CrawlSearchAdd(ctx *gin.Context) {
	domainStr := ctx.Query("domainStr")
	hostStr := ctx.PostForm("host")
	maxDeepLevel, _ := strconv.Atoi(ctx.PostForm("maxDeepInt"))

	cookieValue := ctx.PostForm("cookieValue")
	cookieDomain := ctx.PostForm("cookieDomain")
	cookiePath := ctx.PostForm("cookiePath")

	isReload := ctx.Query("isReload")
	if isReload == "1" {
		if domainStr == "" || maxDeepLevel <= 0 {
			ctx.JSON(http.StatusOK, c.JsonEncode(101, "添加任务失败，参数不正确", nil, 0))
			return
		}
		taskObject := fast_crawl_engine.FastCrawlEngineParams{
			BaseDomain:   domainStr,
			DomainStr:    domainStr,
			MinDeepLevel: 1,
			MaxDeepLevel: maxDeepLevel,
			//Cookies: &fast_crawl_engine.FastCrawlCookies{
			//	Value:  "请使用自己的cookie",
			//	Domain: ".xxx.com",
			//	Path:   "/",
			//},
			//Cookies: nil,
			//Host: "127.0.0.1",
		}
		if cookieValue != "" && cookieDomain != "" && cookiePath != "" {
			taskObject.Cookies = &fast_crawl_engine.FastCrawlCookies{
				Value:  cookieValue,
				Domain: cookieDomain,
				Path:   cookiePath,
			}
		}
		if hostStr != "" {
			taskObject.Host = hostStr
		}
		by, err := json.Marshal(taskObject)
		if err != nil {
			ctx.JSON(http.StatusOK, c.JsonEncode(101, "添加爬虫任务失败，无法开启扫描。"+err.Error(), nil, 0))
			return
		}
		fast_driver.RedisDriver.RPush(fast_driver.RedisWaitCrawlKey, string(by))
		ctx.JSON(http.StatusOK, c.JsonEncode(0, "添加成功", nil, 0))
		return
	}
	ctx.JSON(http.StatusOK, c.JsonEncode(101, "添加任务失败，非正常操作", nil, 0))
	return
}

// 某个域名爬虫结果列表数据
func (c *IndexController) CrawlSearchJson(ctx *gin.Context) {
	domainStr := ctx.Query("domainStr")
	start, _ := strconv.Atoi(ctx.Query("start"))
	r := fast_driver.RedisDriver.LRange(domainStr+fast_driver.RedisDomainCrawlSuffixSymbolKey, int64(start), -1).Val()
	var tmp = make([]struct {
		UrlStr string
		Title  string
	}, 0)
	var rr struct {
		UrlStr string
		Title  string
	}
	for _, v := range r {
		err := json.Unmarshal([]byte(v), &rr)
		if err != nil {
			log.Println(err)
			continue
		}
		tmp = append(tmp, rr)
	}
	ctx.JSON(http.StatusOK, c.JsonEncode(0, gofunc.Md5Encrypt(strings.Join(r, "")), tmp, len(tmp)))
}
