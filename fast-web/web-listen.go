/**
 * Author: haoshuaiwei 
 * Date: 2019-09-18 09:58 
 */

package fast_web

import (
	"context"
	"fast-work/fast-dns-search"
	"fast-work/fast-driver"
	"log"
	"time"
)

var (
	goingBlastMap = map[string]context.CancelFunc{}
)

// 监听待爆破的域名列表
func listenWaitBlastList() {
	log.Println("开始监听待爆破的域名列表...")
	go func() {
		for {
			r := fast_driver.RedisDriver.LPop(fast_driver.RedisWaitBlastKey).Val()
			if len(r) == 0 {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			// 创建上下文对象 添加到正在爆破的map
			go func(goingBlastMap map[string]context.CancelFunc, r string) {
				ctx, cancel := context.WithCancel(context.Background())
				goingBlastMap[r] = cancel
				fast_dns_search.Api(r, ctx)
			}(goingBlastMap, r)
		}
	}()
}

// 监听需要关闭的域名列表
func listenWaitCloaseBlastList() {
	log.Println("开始监听需要关闭的爆破域名列表...")
	go func() {
		for {
			r := fast_driver.RedisDriver.LPop(fast_driver.RedisWaitCloseBlastKey).Val()
			if len(r) == 0 {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			// 关闭正在爆破的域名
			if _, ok := goingBlastMap[r]; !ok {
				log.Println(r, " 关闭失败")
				time.Sleep(time.Millisecond * 100)
				continue
			}
			// 执行关闭 context
			goingBlastMap[r]()
			delete(goingBlastMap, r)
		}
	}()
}
