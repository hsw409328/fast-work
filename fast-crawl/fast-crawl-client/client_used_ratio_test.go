/**
 * Author: haoshuaiwei 
 * Date: 2019-05-23 15:36 
 */

package fast_crawl_client

import (
	"sync"
	"testing"
)

func TestMonitoryClient(t *testing.T) {
	var s sync.WaitGroup
	s.Add(1)
	MonitorClient()
	s.Wait()
}
