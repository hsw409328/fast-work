/**
 * Author: haoshuaiwei
 * Date: 2019-05-15 11:26
 */

package fast_crawl_engine

import (
	"path"
	"strings"
)

func FilterNetWorkRequest(string2 string) bool {
	filterRule := []string{".js", ".png", ".jpg", ".gif", ".flv", ".css", ".woff", ".font"}
	for _, v := range filterRule {
		splitArr := strings.Split(strings.ToLower(path.Ext(string2)), "?")
		if splitArr[0] == v {
			return true
		}
	}
	return false
}
