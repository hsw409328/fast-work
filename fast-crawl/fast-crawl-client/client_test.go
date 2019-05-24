/**
 * Author: haoshuaiwei 
 * Date: 2019-05-15 11:42 
 */

package fast_crawl_client

import (
	"fast-work/fast-crawl/fast-crawl-engine"
	"github.com/benmanns/goworker"
	"testing"
)

func TestClient_Baidu(t *testing.T) {
	goworker.Enqueue(&goworker.Job{
		Queue: "crawl",
		Payload: goworker.Payload{
			Class: "Crawl",
			Args: []interface{}{fast_crawl_engine.FastCrawlEngineParams{
				BaseDomain:   "http://security.jd.com",
				DomainStr:    "http://security.jd.com",
				MinDeepLevel: 1,
				MaxDeepLevel: 4,
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
				Cookies: &fast_crawl_engine.FastCrawlCookies{
					Value:  "请使用自己的淘宝COOKIE",
					Domain: ".taobao.com",
					Path:   "/",
				},
			}},
		},
	})
	Client()
}
