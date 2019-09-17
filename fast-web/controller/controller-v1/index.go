/**
 * Author: haoshuaiwei 
 * Date: 2019-09-17 16:52 
 */

package controller_v1

import (
	"fast-work/fast-web/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IndexController struct {
	controller.BaseController
}

// 功能导航页
func (c *IndexController) Get(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "v1-base-domain/index.html", gin.H{})
}

// 搜索页
func (c *IndexController) DnsSearch(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "v1-base-domain/dns-search.html", gin.H{})
}
