/**
 * Author: haoshuaiwei 
 * Date: 2019-05-15 11:42 
 */

package fast_crawl_client

import (
	"fast-work/fast-crawl/fast-crawl-engine"
	"github.com/benmanns/goworker"
	"log"
	"testing"
	"time"
)

func TestClient_Baidu(t *testing.T) {
	goworker.Enqueue(&goworker.Job{
		Queue: "crawl",
		Payload: goworker.Payload{
			Class: "Crawl",
			Args: []interface{}{fast_crawl_engine.FastCrawlEngineParams{
				BaseDomain:   "http://www.baidu.com",
				DomainStr:    "http://www.baidu.com",
				MinDeepLevel: 1,
				MaxDeepLevel: 2,
				//Cookies: &fast_crawl_engine.FastCrawlCookies{
				//	Value:  "请使用自己的百度cookie",
				//	Domain: ".baidu.com",
				//	Path:   "/",
				//},
				//Cookies: nil,
				//Host: "127.0.0.1",
			}},
		},
	})
	Client()
}

func TestClient_Taobao(t *testing.T) {
	goworker.Enqueue(&goworker.Job{
		Queue: "crawl",
		Payload: goworker.Payload{
			Class: "Crawl",
			Args: []interface{}{fast_crawl_engine.FastCrawlEngineParams{
				BaseDomain:   "https://buyertrade.taobao.com",
				DomainStr:    "https://buyertrade.taobao.com/trade/itemlist/list_bought_items.htm",
				MinDeepLevel: 1,
				MaxDeepLevel: 1,
				//Cookies: &fast_crawl_engine.FastCrawlCookies{
				//	Value:  "请使用自己的淘宝COOKIE",
				//	Domain: ".taobao.com",
				//	Path:   "/",
				//},
			}},
		},
	})
	Client()
}

func TestClient_QQ(t *testing.T) {
	goworker.Enqueue(&goworker.Job{
		Queue: "crawl",
		Payload: goworker.Payload{
			Class: "Crawl",
			Args: []interface{}{fast_crawl_engine.FastCrawlEngineParams{
				BaseDomain:   "https://news.qq.com",
				DomainStr:    "https://news.qq.com",
				MinDeepLevel: 1,
				MaxDeepLevel: 2,
				//Cookies: &fast_crawl_engine.FastCrawlCookies{
				//	Value:  "请使用自己的淘宝COOKIE",
				//	Domain: ".taobao.com",
				//	Path:   "/",
				//},
			}},
		},
	})
	Client()
}

func TestClient_Ability(t *testing.T) {
	//http://testphp.vulnweb.com/
	goworker.Enqueue(&goworker.Job{
		Queue: "crawl",
		Payload: goworker.Payload{
			Class: "Crawl",
			Args: []interface{}{fast_crawl_engine.FastCrawlEngineParams{
				BaseDomain:   "http://testphp.vulnweb.com",
				DomainStr:    "http://testphp.vulnweb.com",
				MinDeepLevel: 1,
				MaxDeepLevel: 4,
			}},
		},
	})
	//Client()
}

func TestClient_51hsw(t *testing.T) {
	//http://testphp.vulnweb.com/
	goworker.Enqueue(&goworker.Job{
		Queue: "crawl",
		Payload: goworker.Payload{
			Class: "Crawl",
			Args: []interface{}{fast_crawl_engine.FastCrawlEngineParams{
				BaseDomain:   "http://www.51hsw.com",
				DomainStr:    "http://www.51hsw.com",
				MinDeepLevel: 1,
				MaxDeepLevel: 4,
			}},
		},
	})
	Client()
}

func TestClient_Coroutine(t *testing.T) {
	//http://testphp.vulnweb.com/
	var b chan bool
	go func() {
		log.Println("启动1")
		Client()
	}()
	go func() {
		time.Sleep(time.Second * 5)
		log.Println("开始执行任务")

		goworker.Enqueue(&goworker.Job{
			Queue: "crawl",
			Payload: goworker.Payload{
				Class: "Crawl",
				Args: []interface{}{fast_crawl_engine.FastCrawlEngineParams{
					BaseDomain:   "http://www.51hsw.com",
					DomainStr:    "http://www.51hsw.com",
					MinDeepLevel: 1,
					MaxDeepLevel: 2,
				}},
			},
		})
	}()
	go func() {
		time.Sleep(time.Second * 15)
		log.Println("开始执行任务 15秋")

		goworker.Enqueue(&goworker.Job{
			Queue: "crawl",
			Payload: goworker.Payload{
				Class: "Crawl",
				Args: []interface{}{fast_crawl_engine.FastCrawlEngineParams{
					BaseDomain:   "http://www.51hsw.com",
					DomainStr:    "http://www.51hsw.com",
					MinDeepLevel: 1,
					MaxDeepLevel: 2,
				}},
			},
		})
	}()
	<-b
}
