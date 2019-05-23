/**
 * Author: haoshuaiwei 
 * Date: 2019-05-14 16:27 
 */

package main

import (
	"github.com/hsw409328/gofunc"
	"log"
)

func main()  {
	log.Println(gofunc.Strpos("//news-bos.cdn.bcebos.com/mvideo/baidu_news_protocol.html","/news.baidu.com"))
	log.Println(gofunc.Strpos("//news-bos.cdn.bcebos.com/mvideo/baidu_news_protocol.html","http://"))
	log.Println(gofunc.Strpos("//news-bos.cdn.bcebos.com/mvideo/baidu_news_protocol.html","https://"))
}