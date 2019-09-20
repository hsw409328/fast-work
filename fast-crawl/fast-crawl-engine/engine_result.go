/**
 * Author: haoshuaiwei 
 * Date: 2019-05-15 11:26 
 */

package fast_crawl_engine

import (
	"encoding/json"
	"fast-work/fast-driver"
	"fast-work/fast-log"
	"fast-work/fast-sys"
	"github.com/benmanns/goworker"
	"github.com/hsw409328/gofunc"
	"log"
	"strings"
)

type FastCrawlResultData struct {
	BaseDomain   string            `基础域名`
	UrlStr       string            `当前URL`
	Title        string            `当前标题`
	Method       string            `请求方式`
	DeepLevel    int               `深度`
	MaxDeepLevel int               `最大深度`
	Cookies      *FastCrawlCookies `cookie`
	Host         string            `服务器地址 可为空`
}

type FastCrawlResult struct {
	result []*FastCrawlResultData
}

func NewFastCrawlResult() *FastCrawlResult {
	return &FastCrawlResult{
		result: make([]*FastCrawlResultData, 0),
	}
}

func (c *FastCrawlResult) Add(result FastCrawlResultData) {
	c.result = append(c.result, &result)
}

func (c *FastCrawlResult) All() []*FastCrawlResultData {
	return c.result
}

func (c *FastCrawlResult) PrintString() string {
	for _, v := range c.result {
		log.Println("【深度】：", v.DeepLevel, "  【请求】：", v.UrlStr, "   ", v.Title)
	}
	return ""
}

func (c *FastCrawlResult) Save() {
	//获取保存的位置
	saveMod, err := fast_sys.GoConfig.GetString("dns", "save_mod")
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range c.result {
		if gofunc.Strpos(saveMod, "redis") {
			by, err := json.Marshal(map[string]string{
				"UrlStr": v.UrlStr,
				"Title":  v.Title,
			})
			if err != nil {
				fast_log.FastLog.Error(err)
			}
			//
			//爬虫结果存储结构使用 具体域名+list+success
			//例如：
			//xx.cc.com-crawl-list-success
			//	{
			//		UrlStr string
			//		Title string
			//	}
			v.BaseDomain = strings.Replace(v.BaseDomain, "http://", "", -1)
			v.BaseDomain = strings.Replace(v.BaseDomain, "https://", "", -1)
			redisCmd := fast_driver.RedisDriver.RPush(v.BaseDomain+"-crawl-list-success", string(by))
			if _, err := redisCmd.Result(); err != nil {
				fast_log.FastLog.Error(err)
			}
		}
		if gofunc.Strpos(saveMod, "mysql") {
			//TODO 暂不支持
		}
	}
}

func (c *FastCrawlResult) SendTask() {
	for _, v := range c.result {
		if (v.DeepLevel + 1) > v.MaxDeepLevel {
			return
		}
		err := goworker.Enqueue(&goworker.Job{
			Queue: "crawl",
			Payload: goworker.Payload{
				Class: "Crawl",
				Args: []interface{}{FastCrawlEngineParams{
					BaseDomain:   v.BaseDomain,
					DomainStr:    v.UrlStr,
					MinDeepLevel: v.DeepLevel + 1,
					MaxDeepLevel: v.MaxDeepLevel,
					Host:         v.Host,
					Cookies:      v.Cookies,
				}},
			},
		})
		if err != nil {
			log.Println(err)
		}
	}

}
