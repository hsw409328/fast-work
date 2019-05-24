/**
 * Author: haoshuaiwei 
 * Date: 2019-05-23 18:02 
 */

package fast_crawl_server

import (
	"errors"
	"fast-work/fast-crawl/fast-crawl-engine"
	"github.com/benmanns/goworker"
)

func AddTask(taskObject fast_crawl_engine.FastCrawlEngineParams) error {
	//taskObject := fast_crawl_engine.FastCrawlEngineParams{
	//	BaseDomain:   "http://news.baidu.com",
	//	DomainStr:    "http://news.baidu.com",
	//	MinDeepLevel: 1,
	//	MaxDeepLevel: 2,
	//	//Cookies: &fast_crawl_engine.FastCrawlCookies{
	//	//	Value:  "请使用自己的百度cookie",
	//	//	Domain: ".baidu.com",
	//	//	Path:   "/",
	//	//},
	//	//Cookies: nil,
	//	//Host: "127.0.0.1",
	//}
	if taskObject.BaseDomain == "" {
		return errors.New("根域名不能为空")
	}
	if taskObject.DomainStr == "" {
		return errors.New("抓取域名不能为空")
	}
	if taskObject.MinDeepLevel <= 0 {
		return errors.New("最小深度必须大于0")
	}
	if taskObject.MaxDeepLevel <= 0 {
		return errors.New("最大深度必须大于1，且大于最小深度值")
	}

	err := goworker.Enqueue(&goworker.Job{
		Queue: "crawl",
		Payload: goworker.Payload{
			Class: "Crawl",
			Args:  []interface{}{taskObject},
		},
	})
	return err
}
