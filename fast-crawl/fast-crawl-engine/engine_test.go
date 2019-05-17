/**
 * Author: haoshuaiwei 
 * Date: 2019-05-15 11:26 
 */

package fast_crawl_engine

import (
	"testing"
)

func TestFastCrawlEngine_Start(t *testing.T) {
	NewFastCrawlEngine(FastCrawlEngineParams{
		BaseDomain:   "http://www.51hsw.com",
		DomainStr:    "http://www.51hsw.com",
		MinDeepLevel: 1,
		MaxDeepLevel: 2,
	}).Start()
}
