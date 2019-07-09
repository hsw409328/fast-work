/**
 * Author: haoshuaiwei
 * Date: 2019-05-15 11:26
 */

package fast_crawl_engine

import (
	"testing"
)

func TestFastCrawlEngine_Start(t *testing.T) {
	var ch chan bool
	NewFastCrawlEngine(FastCrawlEngineParams{
		BaseDomain:   "http://i.baidu.com",
		DomainStr:    "http://i.baidu.com",
		MinDeepLevel: 1,
		MaxDeepLevel: 2,
	}).Start()
	<-ch
}
