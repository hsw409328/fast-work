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
		BaseDomain:   "http://security.jd.com",
		DomainStr:    "http://security.jd.com",
		MinDeepLevel: 1,
		MaxDeepLevel: 1,
	}).Start()
}
