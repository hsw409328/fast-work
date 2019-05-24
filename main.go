/**
 * Author: haoshuaiwei 
 * Date: 2019-05-14 16:27 
 */

package main

import (
	"log"
	"net/url"
)

func main()  {
	urlParse,err := url.Parse("//a.a.com?r=http://a.a.com/asasdf")
	log.Println(err)
	if err!=nil{
		log.Println(err)
	}
	log.Println(urlParse.Host)

	log.Println("http:"+"//www.jd.com/phb/key_737d01ea26d1df7bfe1.html")
}