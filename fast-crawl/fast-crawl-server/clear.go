/**
 * Author: haoshuaiwei 
 * Date: 2019-05-23 18:02 
 */

package fast_crawl_server

import "fast-work/fast-drive"

// 清空指定的bloom-filter 过滤器
func ClearBloomFilter(keyStr string) {
	cmdObject := fast_drive.RedisDriver.Del(keyStr)
	if cmdObject.Err() != nil {
		serverLog.Error(cmdObject.Err().Error())
	}
}
