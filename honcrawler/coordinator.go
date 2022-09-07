package honcrawler

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

/*
	The coordinator manager all the concurrent behaviors.
*/

type coordinator struct {
	TagChannel chan string     // channel for tag map
	HonChannel chan *HonDetail // channel for hon data

	GWaitGroup sync.WaitGroup
	DWaitGroup sync.WaitGroup
	Collector  colly.Collector
}

var Coordinator = coordinator{
	HonChannel: make(chan *HonDetail, 4096),
	GWaitGroup: sync.WaitGroup{},
	DWaitGroup: sync.WaitGroup{},
}

// base collector
var collector = colly.NewCollector(
	colly.UserAgent(UserAgent),
	colly.AllowURLRevisit(),
)

func init() {
	collector.SetRequestTimeout(120 * time.Second)
	// Error Handler
	collector.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error %s: %v\n", r.Request.URL, err)
	})
	// proxy setting
	collector.SetProxy("http://127.0.0.1:55723")
}

func (c *coordinator) Start() {
	c.generate(7101)
	c.consume(4)

	c.GWaitGroup.Wait()
	close(c.HonChannel)
	c.DWaitGroup.Wait()
}

func (c *coordinator) send(hd *HonDetail) {
	c.HonChannel <- hd
}

func (c *coordinator) generate(pCnt int) {
	for i := 1; i <= pCnt; i++ {
		c.GWaitGroup.Add(1)
		go func(i int) {
			infos := GenGallaryInfos(i)
			for _, info := range infos {
				d := GenHonDetails(info)
				c.send(d)
			}
			c.GWaitGroup.Done()
		}(i)
	}

}

func (c *coordinator) consume(cnt int) {
	for i := 0; i < cnt; i++ {
		c.DWaitGroup.Add(1)
		go func(i int) {
			for hon := range c.HonChannel {
				Download(hon)
			}
			c.DWaitGroup.Done()
		}(i)
	}
}

func (c *coordinator) Try() {
	c.generate(1)
	c.consume(4)

	c.GWaitGroup.Wait()
	close(c.HonChannel)
	c.DWaitGroup.Wait()
}
