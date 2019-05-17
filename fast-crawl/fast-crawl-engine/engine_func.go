/**
 * Author: haoshuaiwei 
 * Date: 2019-05-15 11:26 
 */

package fast_crawl_engine

import (
	"gofunc"
	"strings"
)

func FilterNetWorkRequest(string2 string) bool {
	filterRule := []string{".js", ".png", ".jpg", ".gif", ".flv", ".css"}
	for _, v := range filterRule {
		if gofunc.Strpos(strings.ToLower(string2), v) {
			return true
		}
	}
	return false
}
