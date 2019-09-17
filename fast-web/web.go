/**
 * Author: haoshuaiwei 
 * Date: 2019-05-14 16:27 
 */

package fast_web

import (
	"fast-work/fast-sys"
	"fast-work/fast-web/controller"
	"fast-work/fast-web/controller/controller-v1"
	"github.com/gin-gonic/gin"
)

func Web() {

	// 启动WEB管理
	en := gin.Default()
	en.Delims("[[[", "]]]")
	en.Static("/layui", "./fast-web/static/layui")

	en.LoadHTMLGlob("./fast-web/views/**/*")

	en.GET("/", (&controller.IndexController{}).Get)

	// v1版本
	v1 := en.Group("/v1")
	{
		v1.GET("/", (&controller_v1.IndexController{}).Get)
		v1.GET("/json", (&controller_v1.IndexController{}).GetJson)
		v1.GET("/dns/search", (&controller_v1.IndexController{}).DnsSearch)
	}

	webRunConfig, _ := fast_sys.GoConfig.GetString("web", "run_config")
	en.Run(webRunConfig)
}
