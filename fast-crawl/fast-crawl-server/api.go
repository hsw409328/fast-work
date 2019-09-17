/**
 * Author: haoshuaiwei 
 * Date: 2019-05-23 15:31 
 */

package fast_crawl_server

import (
	"fast-work/fast-sys"
	"github.com/hsw409328/gofunc"
	"github.com/hsw409328/gofunc/go_hlog"
	"time"
)

var (
	serverLog         *go_hlog.Logger
	logDateSymbol     string
	goWorkerQueueName = "crawl"
	goWorkerNamespace = "fast-crawl:"
	redisAddress, _   = fast_sys.GoConfig.GetString("redis", "host")
	redisPort, _      = fast_sys.GoConfig.GetString("redis", "port")
)

func initLog() {
	logDateSymbol = gofunc.CurrentDate()
	serverLog = go_hlog.GetInstance(gofunc.GetCurrentPath() + "/log-file/" + gofunc.CurrentDate() + ".log")
	go func() {
		// 监听自动切割系统运行日志
		for {
			if logDateSymbol != gofunc.CurrentDate() {
				logDateSymbol = gofunc.CurrentDate()
				serverLog = go_hlog.GetInstance(gofunc.GetCurrentPath() + "/log-file/" + gofunc.CurrentDate() + ".log")
			}
			time.Sleep(time.Second)
		}
	}()
}

func init() {
	initLog()
}

// 统一对外提供使用 API方法调用
func Api() {
	//管理在线客户端
	//添加待扫描任务
	//下发客户端指令 --stop --restart
	//清空bloom_filter的key
}