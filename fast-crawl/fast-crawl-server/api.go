/**
 * Author: haoshuaiwei 
 * Date: 2019-05-23 15:31 
 */

package fast_crawl_server

import (
	"fast-work/fast-crawl/fast-crawl-client"
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
	go Api("start")
}

// 统一对外提供使用 API方法调用
func Api(cmdStr string) interface{} {
	switch cmdStr {
	case "start":
		// 启动爬虫客户端扫描队列
		fast_crawl_client.Client()
	case "stop":
	case "clear_bloom_filter":
		ClearBloomFilter("fastCrawlKey")
	case "client_list":
		// 启动
		return (&ManagerClient{}).Read()
	default:
		break
	}
	return nil
}
