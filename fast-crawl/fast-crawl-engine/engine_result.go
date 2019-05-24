/**
 * Author: haoshuaiwei 
 * Date: 2019-05-15 11:26 
 */

package fast_crawl_engine

import (
	"github.com/benmanns/goworker"
	"log"
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
		log.Println("【深度】：", v.DeepLevel, "  【请求】：", v.UrlStr)
	}
	return ""
}

func (c *FastCrawlResult) SendTask() {
	for _, v := range c.result {
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
