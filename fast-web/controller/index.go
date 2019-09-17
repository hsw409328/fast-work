/**
 * Author: haoshuaiwei 
 * Date: 2019-09-17 10:56 
 */

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IndexController struct {
	BaseController
}

// 首页
func (c *IndexController) Get(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index/index.html", gin.H{})
}
