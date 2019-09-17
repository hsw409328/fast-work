/**
 * Author: haoshuaiwei 
 * Date: 2019-09-17 16:52 
 */

package controller_v1

import (
	"fast-work/fast-driver"
	"fast-work/fast-web/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	dnsBlastKey = "dns-search-list"
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
	r := fast_driver.RedisDriver.LRange(dnsBlastKey, 0, -1).Val()
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
	r := fast_driver.RedisDriver.LRange(dnsBlastKey, 0, -1).Val()
	for _, v := range r {
		if domainStr == v {
			return true
		}
	}
	return false
}

// 搜索页
func (c *IndexController) DnsSearch(ctx *gin.Context) {
	domainStr := ctx.Query("domainStr")
	err := c.CheckDnsBlastDomainIsExist(domainStr)
	if !err {
		fast_driver.RedisDriver.RPush(dnsBlastKey, domainStr)
	}
	ctx.HTML(http.StatusOK, "v1-base-domain/dns-search.html", gin.H{"baseDomain": domainStr})
}
