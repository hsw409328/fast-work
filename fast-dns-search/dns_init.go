/**
 * Author: haoshuaiwei 
 * Date: 2019-05-14 17:02 
 */

package fast_dns_search

import (
	"fast-work/fast-log"
	"github.com/hsw409328/gofunc"
)

func init() {
	domainDict, err = gofunc.ReadLinesForFile("dict/list.txt")
	if err != nil {
		fast_log.FastLog.Error(err)
	}
}
