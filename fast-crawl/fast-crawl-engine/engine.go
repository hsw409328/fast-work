/**
 * Author: haoshuaiwei 
 * Date: 2019-05-15 11:26 
 */

package fast_crawl_engine

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/hsw409328/gofunc"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	jsStr = `try {
            	window.alert = function(msg) {};
           		window.confirm = function(msg) {
               		return false
           		};
           		window.prompt = function(text, defaultText) {
               		return false
           		};
           		window.close = function() {
               		return false
           		};
           		window.history.back = function(args) {
               		console.log(args)
           		};
           		window.history.forward = function(args) {
               		console.log(args)
           		};
				var f = function(){
					var eles = document.getElementsByTagName('*');
					for (x in eles) {
						elm = eles[x];
						if(typeof elm.click !== "undefined"){elm.click();}
					}
				}
				f();
       		} catch (err) {
           		console.log(err)
       		}
		`
)

type FastCrawlCookies struct {
	Value  string
	Domain string
	Path   string
}

type FastCrawlEngineParams struct {
	BaseDomain   string            `根域名`
	DomainStr    string            `渲染抓取的域名`
	Cookies      *FastCrawlCookies `cookie`
	Host         string            `服务器地址 可为空`
	MinDeepLevel int               `基础深度`
	MaxDeepLevel int               `最大深度`
}

type FastCrawlEngine struct {
	params *FastCrawlEngineParams
}

// 外部实例
func NewFastCrawlEngine(params FastCrawlEngineParams) *FastCrawlEngine {
	return &FastCrawlEngine{
		params: &params,
	}
}

// 启动扫描
func (c *FastCrawlEngine) Start() {
	if c.params.MinDeepLevel > c.params.MaxDeepLevel {
		return
	}

	resultMap := sync.Map{}
	resultObject := NewFastCrawlResult()

	var htmlStr string
	var jsInterface interface{}
	c.initRender(&htmlStr, &jsInterface, &resultMap)

	//解析需要的URL连接
	var parseHtml = func(htmlStr string) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
		if err != nil {
			log.Println(err)
		}
		title := doc.Find("title").Text()
		doc.Find("a").Each(func(i int, selection *goquery.Selection) {
			urlStr, _ := selection.Attr("href")
			resultMap.Store(urlStr, "get")
		})
		doc.Find("iframe").Each(func(i int, selection *goquery.Selection) {
			urlStr, _ := selection.Attr("src")
			resultMap.Store(urlStr, "get")
		})
		doc.Find("form").Each(func(i int, selection *goquery.Selection) {
			urlStr, _ := selection.Attr("action")
			methodStr, _ := selection.Attr("method")
			if methodStr == "" || strings.ToLower(methodStr) == "get" {
				resultMap.Store(urlStr, "get")
			} else {
				resultMap.Store(urlStr, "post")
			}
		})

		// 添加的本次获取的结果集
		resultMap.Range(func(key, value interface{}) bool {
			k := gofunc.InterfaceToString(key)
			//替换http:// 和 https:// 防止判断完整性错误，以及 域名的子域名，例如：i.xxx.com和a.i.xxx.com
			tmpBaseDomain := strings.Replace(c.params.BaseDomain, "http://", "/", -1)
			tmpBaseDomain = strings.Replace(tmpBaseDomain, "https://", "/", -1)
			// 判断连接是否完整
			if !gofunc.Strpos(k, tmpBaseDomain) {
				//判断域名最后一个是否有 "/" and 判断原始连接第一个是否有 "/" 进行连接
				if k == "" {
					k = "/"
				}
				k = c.params.BaseDomain + gofunc.ConnectFirstWord(k, "/")
			}

			resultObject.Add(FastCrawlResultData{
				BaseDomain:   c.params.BaseDomain,
				UrlStr:       k,
				Method:       gofunc.InterfaceToString(value),
				Title:        title,
				DeepLevel:    c.params.MinDeepLevel,
				MaxDeepLevel: c.params.MaxDeepLevel,
				Host:         c.params.Host,
				Cookies:      c.params.Cookies,
			})
			return true
		})
	}
	//调用解析，添加到任务列表
	parseHtml(htmlStr)

	//时时查看结果
	resultObject.PrintString()

	//添加每一轮的结果到消息队列中
	resultObject.SendTask()
}

// 渲染引擎
func (c *FastCrawlEngine) initRender(htmlStr *string, jsInterface *interface{}, r *sync.Map) error {
	var browserOptions []chromedp.ContextOption
	browserOptions = append(browserOptions,
		// 拦截网络请求
		chromedp.WithDebugf(func(s string, i ...interface{}) {
			go func() {
				for _, elem := range i {
					var msg cdproto.Message
					var msgIn struct {
						SessionId string `json:"sessionId"`
						Message   string `json:"message"`
					}
					var msgLast cdproto.Message
					// The CDP messages are sent as strings so we need to convert them back
					err := json.Unmarshal([]byte(fmt.Sprintf("%s", elem)), &msg)
					if err != nil {
						continue
					}
					err = json.Unmarshal(msg.Params, &msgIn)
					if err != nil {
						continue
					}
					err = json.Unmarshal([]byte(msgIn.Message), &msgLast)
					//log.Println(string(msgIn.Message))
					// 拦截请求
					// Network.requestWillBeSent {"requestId":"","loaderId":"","documentURL":"","request":{"url":"http://xxx.xx.com/","method":"GET"}}
					var BeSent struct {
						Request struct {
							Url    string
							Method string
						}
					}
					// Page.navigatedWithinDocument {"frameId":"","url":"http://xxx.xx.com/#/notice/196"}
					// Page.windowOpen {"url":"http://xxx.xx.com/#/user/xxx"}
					// Page.frameScheduledNavigation {"url":"https://passport.jd.com/uc/login?ReturnUrl=http://xxx.xx.com/#/"}
					var CommonEvent struct {
						Url string
					}
					switch msgLast.Method.String() {
					case "Network.requestWillBeSent":
						by, _ := msgLast.Params.MarshalJSON()
						json.Unmarshal(by, &BeSent)
						if !FilterNetWorkRequest(BeSent.Request.Url) {
							r.Store(BeSent.Request.Url, strings.ToLower(BeSent.Request.Method))
						}
						break
					case "Page.navigatedWithinDocument":
						by, _ := msgLast.Params.MarshalJSON()
						json.Unmarshal(by, &CommonEvent)
						if !FilterNetWorkRequest(CommonEvent.Url) {
							r.Store(CommonEvent.Url, "get")
						}
						break
					case "Page.windowOpen":
						by, _ := msgLast.Params.MarshalJSON()
						json.Unmarshal(by, &CommonEvent)
						if !FilterNetWorkRequest(CommonEvent.Url) {
							r.Store(CommonEvent.Url, "get")
						}
						break
					case "Page.frameScheduledNavigation":
						by, _ := msgLast.Params.MarshalJSON()
						json.Unmarshal(by, &CommonEvent)
						if !FilterNetWorkRequest(CommonEvent.Url) {
							r.Store(CommonEvent.Url, "get")
						}
						break
					}
				}
			}()
		}),
		chromedp.WithBrowserOption(
			chromedp.WithDialTimeout(time.Second*10),
		),
	)
	ctx, cancel := chromedp.NewContext(context.Background(), browserOptions...)
	defer cancel()

	// 设置host
	networkHeaders := network.Headers{}
	if c.params.Host != "" {
		networkHeaders["Host"] = c.params.Host
	}
	// 设置cookie
	networkCookies := []*network.CookieParam{}
	if c.params.Cookies != nil {
		// 拆解cookie 先用";" 再用"=" 需要去空
		tmpCookieList := strings.Split(c.params.Cookies.Value, ";")
		for _, v := range tmpCookieList {
			if v == "" {
				continue
			}
			tmpKeyValue := strings.Split(v, "=")
			networkCookies = append(networkCookies, &network.CookieParam{
				Name:   strings.Replace(tmpKeyValue[0], " ", "", -1),
				Value:  strings.Replace(tmpKeyValue[1], " ", "", -1),
				Domain: c.params.Cookies.Domain,
				Path:   c.params.Cookies.Path,
			})
		}
	}

	// 正式启动引擎
	err := chromedp.Run(ctx,
		chromedp.Tasks{
			network.Enable(),
			network.SetCookies(networkCookies),
			network.SetExtraHTTPHeaders(networkHeaders),
			chromedp.Navigate(c.params.DomainStr),
			chromedp.WaitReady("body", chromedp.ByQuery),
			//用以下的方法有BUG，存在未找到标签的情况下无限循环
			//chromedp.AttributesAll("a", &attributes, chromedp.ByQueryAll),
			//chromedp.AttributesAll("iframe", &attributes, chromedp.ByQueryAll),
			//chromedp.AttributesAll("form", &attributes, chromedp.ByQueryAll),
			// 获取页面所有html
			chromedp.OuterHTML("html", htmlStr, chromedp.ByQuery),
			chromedp.EvaluateAsDevTools(jsStr, jsInterface),
		},
	)

	return err
}
