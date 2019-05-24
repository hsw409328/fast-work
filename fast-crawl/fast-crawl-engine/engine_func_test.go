/**
 * Author: haoshuaiwei
 * Date: 2019-05-15 11:26
 */

package fast_crawl_engine

import (
	"testing"
)

func TestFilterNetWorkRequest(t *testing.T) {
	a := "https://news.qq.com/ext2020/apub/json/prevent.new.json"
	data := FilterNetWorkRequest(a)
	expect := false
	if data != expect {
		t.Error("func err")
	}
}
