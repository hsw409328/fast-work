/**
 * Author: haoshuaiwei 
 * Date: 2019-05-14 16:47 
 */

package fast_dns_search

import (
	"encoding/json"
	"fast-work/fast-driver"
	"fast-work/fast-log"
	"fast-work/fast-sys"
	"github.com/hsw409328/gofunc"
	"github.com/hsw409328/gofunc/go_hlog"
	"log"
)

type DnsResultData struct {
	BaseDomain string
	Domain     string
	Ip         []string
	Host       string
	IsOpen     bool   `是否可以打开`
	Type       string `结果类型： A\TXT\CNAME\IP`
}

type DnsResult struct {
	saveLog *go_hlog.Logger
}

func NewDnsResult(logger *go_hlog.Logger) *DnsResult {
	return &DnsResult{
		saveLog: logger,
	}
}

func (c *DnsResult) Save(data DnsResultData) {
	//获取保存的位置
	saveMod, err := fast_sys.GoConfig.GetString("dns", "save_mod")
	if err != nil {
		log.Fatal(err)
	}
	if gofunc.Strpos(saveMod, "local") {
		c.saveLog.Info(data)
	}
	if gofunc.Strpos(saveMod, "redis") {
		by, err := json.Marshal(data)
		if err != nil {
			fast_log.FastLog.Error(err)
		}
		redisCmd := fast_driver.RedisDriver.HSet("dns_"+data.BaseDomain, data.Domain, string(by))
		if _, err := redisCmd.Result(); err != nil {
			fast_log.FastLog.Error(err)
		}
	}
	if gofunc.Strpos(saveMod, "mysql") {
		//TODO 暂不支持
	}
}
