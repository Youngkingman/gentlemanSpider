package honcrawler

import (
	"fmt"
	"sync"
	"time"

	"github.com/Youngkingman/gentlemanSpider/settings"
	colly "github.com/gocolly/colly/v2"
)

/*
	Some constant for the spider
*/
const Host = `https://www.wnacg.com`
const GallaryUrl string = Host + `/albums-index-page-%d.html`
const UserAgent string = `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:104.0) Gecko/20100101 Firefox/104.0`
const ImgsPerPage int = 12

/*
	The coordinator manager all the concurrent behaviors.
*/

type coordinator struct {
	tagChannel chan string     // channel for tag map
	honChannel chan *HonDetail // channel for hon data

	gWaitGroup sync.WaitGroup
	dWaitGroup sync.WaitGroup

	tagMap map[string]bool
	mutex  sync.Mutex
}

var Coordinator = coordinator{
	honChannel: make(chan *HonDetail, 4096),
	tagChannel: make(chan string, 4096),
	gWaitGroup: sync.WaitGroup{},
	dWaitGroup: sync.WaitGroup{},
	tagMap:     make(map[string]bool),
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
	collector.SetProxy(settings.CrawlerSetting.ProxyHost)
}

func (c *coordinator) sendHon(hd *HonDetail) {
	c.honChannel <- hd
}

func (c *coordinator) sendTag(tag string) {
	c.tagChannel <- tag
}

func (c *coordinator) generateHon(pSt int, pEnd int) {
	for i := pSt; i <= pEnd; i++ {
		c.gWaitGroup.Add(1)
		go func(i int) {
			infos := GenGallaryInfos(i)
			for _, info := range infos {
				d := GenHonDetails(info)
				for _, t := range d.Tags {
					c.sendTag(t)
				}
				c.sendHon(d)
			}
			c.gWaitGroup.Done()
		}(i)
	}

}

func (c *coordinator) consumeHon(cnt int) {
	for i := 0; i < cnt; i++ {
		c.dWaitGroup.Add(1)
		go func() {
			for hon := range c.honChannel {
				Download(hon)
			}
			c.dWaitGroup.Done()
		}()
	}
}

func (c *coordinator) comsumeTag(cnt int) {
	for i := 0; i < cnt; i++ {
		c.dWaitGroup.Add(1)
		go func() {
			for tag := range c.tagChannel {
				c.mutex.Lock()
				if ok := c.tagMap[tag]; !ok {
					c.tagMap[tag] = true
					SaveTag(tag)
				}
				c.mutex.Unlock()
			}
		}()
	}
}

func (c *coordinator) Start() {
	c.generateHon(
		settings.CrawlerSetting.PageStart,
		settings.CrawlerSetting.PageEnd,
	)
	c.consumeHon(settings.CrawlerSetting.HonConsumerCount)
	c.comsumeTag(settings.CrawlerSetting.TagConsumerCount)

	c.gWaitGroup.Wait()
	close(c.honChannel)
	close(c.tagChannel)
	c.dWaitGroup.Wait()
}
