/**
 * Author: haoshuaiwei 
 * Date: 2019-09-17 10:46 
 */

package controller

import "github.com/gin-gonic/gin"

type BaseController struct {
}

func (c *BaseController) JsonEncode(code int, msg string, data interface{}, count int) gin.H {
	return gin.H{
		"code":  code,
		"msg":   msg,
		"data":  data,
		"count": count,
	}
}
