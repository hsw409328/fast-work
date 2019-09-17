/**
 * Author: haoshuaiwei 
 * Date: 2019-05-23 18:01 
 */

package fast_crawl_server

import (
	"fast-work/fast-driver"
)

const (
	ClientSymbol = "ClientInfo"
)

type ManagerClient struct {
}

func (c *ManagerClient) Read() map[string]string {
	return fast_driver.RedisDriver.HGetAll(ClientSymbol).Val()
}

func (c *ManagerClient) Page() {
	clientMap := c.Read()
	for k, v := range clientMap {
		serverLog.Info(k, v)
	}
}
