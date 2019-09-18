/**
 * Author: haoshuaiwei 
 * Date: 2019-09-17 16:52 
 */

package controller_v1

import (
	"encoding/json"
	"fast-work/fast-dns-search"
	"fast-work/fast-driver"
	"fast-work/fast-web/controller"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
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
	err := c.CheckDnsBlastDomainIsExist(domainStr)
	if !err {
		fast_driver.RedisDriver.RPush(fast_driver.RedisDnsBlastKey, domainStr)
	}
	c.researchDns(domainStr)
	ctx.HTML(http.StatusOK, "v1-base-domain/dns-search.html", gin.H{"baseDomain": domainStr})
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
	ctx.JSON(http.StatusOK, c.JsonEncode(0, "success", tmp, len(tmp)))
}
