/**
 * Author: haoshuaiwei 
 * Date: 2019-05-14 16:27 
 */

package main

import (
	"log"
	"runtime"
)

func main()  {
	log.Println(runtime.NumCPU())
}