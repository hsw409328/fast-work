/**
 * Author: haoshuaiwei 
 * Date: 2019-05-14 16:31 
 */

package fast_dns_search

import (
	"testing"
)

func TestDnsBlast_Run(t *testing.T) {
	//object := NewDnsBlast("baidu.com")
	//object.A = true
	//object.Run()
}

func TestNewDnsBlast(t *testing.T) {
	//go func() {
	//	object := NewDnsBlast("baidu.com")
	//	object.A = true
	//	object.Run()
	//}()
	//go func() {
	//	object1 := NewDnsBlast("weibo.com")
	//	object1.A = true
	//	object1.Run()
	//}()
	//time.Sleep(time.Second * 1)
}

func TestDnsEle(t *testing.T) {
	object := NewDnsBlast("ele.me")
	object.A = true
	object.Run()
}

func TestDnsQq(t *testing.T)  {
	object := NewDnsBlast("tencent.com")
	object.A = true
	object.Run()
}