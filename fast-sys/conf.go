/**
 * Author: haoshuaiwei 
 * Date: 2019-05-14 16:44 
 */

package fast_sys

import (
	"github.com/hsw409328/gofunc"
	"github.com/hsw409328/gofunc/go_config"
)

var (
	GoConfig = go_config.NewReadConfigLib(gofunc.GetCurrentPath() + "/fast.conf")
)
