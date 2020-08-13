package main

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-go-sdk"
	"github.com/crawlab-team/crawlab-go-sdk/entity"
	"github.com/gocolly/colly/v2"
	"runtime/debug"
)

func main() {
	startUrl := "https://segmentfault.com/search?q=crawlab"

	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36"),
	)

	c.OnHTML(".search-result > .widget-blog", func(e *colly.HTMLElement) {
		item := entity.Item{}
		item["title"] = e.ChildText("h2.h4 > a")
		item["url"] = e.ChildAttr("h2.h4 > a", "href")
		fmt.Println(item)
		if err := crawlab.SaveItem(item); err != nil {
			log.Errorf("save item error: " + err.Error())
			debug.PrintStack()
			return
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(fmt.Sprintf("Visiting %s", r.URL.String()))
	})

	if err := c.Visit(startUrl); err != nil {
		log.Errorf("visit error: " + err.Error())
		debug.PrintStack()
		panic(fmt.Sprintf("Unable to visit %s", startUrl))
	}

	c.Wait()
}
