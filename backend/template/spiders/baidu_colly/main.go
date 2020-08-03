package main

import (
	"fmt"
	"github.com/apex/log"
	"github.com/gocolly/colly/v2"
	"runtime/debug"
)

type BaiduItem struct {
	title string
	url   string
}

func main() {
	startUrl := "https://www.baidu.com/s?wd=crawlab"

	c := colly.NewCollector(
		colly.AllowedDomains("www.baidu.com"),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36"),
	)

	c.OnHTML("#content_left > .c-container", func(e *colly.HTMLElement) {
		item := BaiduItem{
			title: e.ChildText("h3.t > a"),
			url:   e.ChildAttr("h3.t > a", "href"),
		}
		fmt.Println(item)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Debugf(fmt.Sprintf("Visiting %s", r.URL.String()))
	})

	if err := c.Visit(startUrl); err != nil {
		log.Errorf("visit error: " + err.Error())
		debug.PrintStack()
		panic(fmt.Sprintf("Unable to visit %s", startUrl))
	}

	c.Wait()
}
