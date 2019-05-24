/**
 * Author: haoshuaiwei 
 * Date: 2019-05-15 11:13 
 */

package fast_crawl_client

import (
	"encoding/json"
	"fast-work/fast-crawl/fast-crawl-engine"
	"fast-work/fast-sys"
	"github.com/benmanns/goworker"
	"github.com/hsw409328/gofunc"
	"github.com/hsw409328/gofunc/go_hlog"
	"log"
	"runtime"
	"time"
)

var (
	clientLog         *go_hlog.Logger
	logDateSymbol     string
	goWorkerQueueName = "crawl"
	goWorkerNamespace = "fast-crawl:"
	redisAddress, _   = fast_sys.GoConfig.GetString("redis", "host")
	redisPort, _      = fast_sys.GoConfig.GetString("redis", "port")
)

func initLog() {
	logDateSymbol = gofunc.CurrentDate()
	clientLog = go_hlog.GetInstance(gofunc.GetCurrentPath() + "/log-file/" + gofunc.CurrentDate() + ".log")
	go func() {
		// 监听自动切割系统运行日志
		for {
			if logDateSymbol != gofunc.CurrentDate() {
				logDateSymbol = gofunc.CurrentDate()
				clientLog = go_hlog.GetInstance(gofunc.GetCurrentPath() + "/log-file/" + gofunc.CurrentDate() + ".log")
			}
			time.Sleep(time.Second)
		}
	}()
}

func init() {
	initLog()
	settings := goworker.WorkerSettings{
		URI:            "redis://" + redisAddress + ":" + redisPort + "/",
		Connections:    100,
		Queues:         []string{goWorkerQueueName},
		UseNumber:      true,
		ExitOnComplete: false,
		Concurrency:    runtime.NumCPU(),
		Namespace:      goWorkerNamespace,
		Interval:       5.0,
	}
	goworker.SetSettings(settings)
	goworker.Register("Crawl", crawl)
}

func crawl(queue string, args ...interface{}) error {
	log.Println("队列名称【", queue, "】--------开始执行---------")
	var fastCrawlParams fast_crawl_engine.FastCrawlEngineParams
	for _, v := range args {
		by, err := json.Marshal(v)
		if err != nil {
			clientLog.Error(err)
			return err
		}
		json.Unmarshal(by, &fastCrawlParams)
		fast_crawl_engine.NewFastCrawlEngine(fastCrawlParams).Start()
	}
	defer log.Println("队列名称【", queue, "】执行完成")
	return nil
}

func Client() {
	if err := goworker.Work(); err != nil {
		log.Fatal(err)
	}
}
