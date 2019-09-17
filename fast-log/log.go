/**
 * Author: haoshuaiwei
 * Date: 2019-05-14 16:40
 */

package fast_log

import (
	"github.com/hsw409328/gofunc"
	"github.com/hsw409328/gofunc/go_hlog"
	"time"
)

var (
	FastLog       *go_hlog.Logger
	logDateSymbol string
)

func init() {
	logDateSymbol = gofunc.CurrentDate()
	FastLog = go_hlog.GetInstance(gofunc.GetCurrentPath() + "/log-file/" + gofunc.CurrentDate() + ".log")
	go func() {
		// 监听自动切割系统运行日志
		for {
			if logDateSymbol != gofunc.CurrentDate() {
				logDateSymbol = gofunc.CurrentDate()
				FastLog = go_hlog.GetInstance(gofunc.GetCurrentPath() + "/log-file/" + gofunc.CurrentDate() + ".log")
			}
			time.Sleep(time.Second)
		}
	}()
}
