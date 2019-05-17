/**
 * Author: haoshuaiwei 
 * Date: 2019-05-14 16:28 
 */

package fast_dns_search

import (
	"fast-work/fast-log"
	"fast-work/fast-sys"
	"fmt"
	"github.com/hsw409328/gofunc"
	"github.com/hsw409328/gofunc/go_hlog"
	"github.com/hsw409328/gofunc/go_pool"
	"net"
	"runtime"
	"strconv"
	"time"
)

var (
	blastCpuMultipleNumber, _ = fast_sys.GoConfig.GetString("dns", "blast_cpu_multiple")
)

type DnsBlast struct {
	A          bool `是否获取A记录`
	Txt        bool `是否获取txt记录`
	Cname      bool `是否获取cname记录`
	Ip         bool `是否读取IP记录`
	LogFile    *go_hlog.Logger
	LogErrFile *go_hlog.Logger
	Domain     string
}

// 设置要爆破的域名
func NewDnsBlast(baseDomain string) *DnsBlast {
	return &DnsBlast{
		Domain: baseDomain,
	}
}

// 读取A记录
func (c *DnsBlast) ReadHostARecord(domain string) {
	fmt.Println(domain)
	time.Sleep(time.Second)
	if c.A {
		ipString, err := net.LookupHost(domain)
		if err != nil {
			c.LogErrFile.Error(err)
			return
		}
		NewDnsResult(c.LogFile).Save(DnsResultData{
			BaseDomain: c.Domain,
			Domain:     domain,
			Ip:         ipString,
			IsOpen:     true,
			Type:       "A",
		})
	}
}

// 读取TXT记录
func (c *DnsBlast) ReadHostTxtRecord(domain string) {
	if c.Txt {
		ipString, err := net.LookupTXT(domain)
		if err != nil {
			c.LogErrFile.Error(err)
			return
		}
		NewDnsResult(c.LogFile).Save(DnsResultData{
			BaseDomain: c.Domain,
			Domain:     domain,
			Ip:         ipString,
			IsOpen:     true,
			Type:       "TXT",
		})
	}
}

// 读取CNAME记录
func (c *DnsBlast) ReadHostCnameRecord(domain string) {
	if c.Cname {
		cnameString, err := net.LookupCNAME(domain)
		if err != nil {
			c.LogErrFile.Error(err)
			return
		}
		NewDnsResult(c.LogFile).Save(DnsResultData{
			BaseDomain: c.Domain,
			Domain:     domain,
			Host:       cnameString,
			IsOpen:     true,
			Type:       "CNAME",
		})
	}
}

// 读取IP记录
func (c *DnsBlast) ReadHostIpRecord(domain string) {
	if c.Ip {
		ipString, err := net.LookupHost(domain)
		if err != nil {
			c.LogErrFile.Error(err)
			return
		}
		NewDnsResult(c.LogFile).Save(DnsResultData{
			BaseDomain: c.Domain,
			Domain:     domain,
			Ip:         ipString,
			IsOpen:     true,
			Type:       "IP",
		})
	}
}

func (c *DnsBlast) start(val interface{}) {
	v := gofunc.InterfaceToString(val)
	c.ReadHostARecord(v)
	c.ReadHostCnameRecord(v)
	c.ReadHostIpRecord(v)
	c.ReadHostTxtRecord(v)
}

// 生成爆破任务
func (c *DnsBlast) createTask() {
	// 判断是否为域名
	if !gofunc.IsDomain(c.Domain) {
		fast_log.FastLog.Error(c.Domain, "非合法域名")
		return
	}
	c.LogFile = go_hlog.GetInstance(gofunc.GetCurrentPath() + "/cache/" + c.Domain + ".log")
	c.LogErrFile = go_hlog.GetInstance(gofunc.GetCurrentPath() + "/cache/" + c.Domain + ".err" + ".log")

	// 产生任务
	tmpCpu, _ := strconv.Atoi(blastCpuMultipleNumber)
	g := go_pool.NewGoPool(runtime.NumCPU()*tmpCpu, c.start)
	go func() {
		for i := 0; i < len(domainDict); i++ {
			g.Push(domainDict[i] + "." + c.Domain)
		}
		g.Close()
	}()
	g.Run()
}

// 执行
func (c *DnsBlast) Run() {
	c.createTask()
}
