/**
 * Author: haoshuaiwei 
 * Date: 2019-07-12 11:58 
 */

package fast_dns_search

// 统一对外提供使用 API方法调用
func Api(domainStr string) {
	object := NewDnsBlast(domainStr)
	object.A = true
	object.Run()
}
