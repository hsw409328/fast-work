/**
 * Author: haoshuaiwei
 * Date: 2019-05-14 16:31
 */

package fast_dns_search

import (
	"context"
	"testing"
	"time"
)

func TestDnsBlast_Run(t *testing.T) {
	//object := NewDnsBlast("baidu.com")
	//object.A = true
	//object.Run()
}

func TestNewDnsBlast(t *testing.T) {
	// 管理子协程
	m := map[string]context.CancelFunc{}
	go func(m map[string]context.CancelFunc) {
		ctx, cancel := context.WithCancel(context.Background())
		m["test1"] = cancel
		object := NewDnsBlast("baidu.com")
		object.A = true
		object.Run(ctx)
	}(m)
	go func(m map[string]context.CancelFunc) {
		ctx, cancel := context.WithCancel(context.Background())
		m["test2"] = cancel
		object1 := NewDnsBlast("weibo.com")
		object1.A = true
		object1.Run(ctx)
	}(m)
	go func() {
		time.Sleep(time.Second * 1)
		m["test1"]()
		time.Sleep(time.Second * 1)
		m["test2"]()
	}()
	time.Sleep(time.Second * 5)
}

func TestDnsEle(t *testing.T) {
	object := NewDnsBlast("ele.me")
	object.A = true
	ctx, cancel := context.WithCancel(context.Background())
	go func(cancel context.CancelFunc) {
		time.Sleep(time.Second * 3)
		cancel()
	}(cancel)
	object.Run(ctx)
}

func TestDnsQq(t *testing.T) {
	object := NewDnsBlast("tencent.com")
	object.A = true
	ctx, cancel := context.WithCancel(context.Background())
	object.Run(ctx)
	cancel()
}
