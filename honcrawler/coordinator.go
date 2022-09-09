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

	tagSet set
	mutex  sync.Mutex
}

var Coordinator = coordinator{
	honChannel: make(chan *HonDetail, settings.CrawlerSetting.HonBuffer),
	tagChannel: make(chan string, settings.CrawlerSetting.TagBuffer),
	gWaitGroup: sync.WaitGroup{},
	dWaitGroup: sync.WaitGroup{},
	// tagSet:     make(map[string]struct{}),
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
	if settings.CrawlerSetting.EnableProxy {
		collector.SetProxy(settings.CrawlerSetting.ProxyHost)
	}
}

func (c *coordinator) sendHon(hd *HonDetail) {
	c.honChannel <- hd
}

func (c *coordinator) sendTag(tag string) {
	c.tagChannel <- tag
}

func (c *coordinator) generateHon(pSt int, pEnd int) {
	ch := make(chan struct{}, settings.CrawlerSetting.HonConsumerCount*3/2)
	for i := pSt; i <= pEnd; i++ {
		ch <- struct{}{}
		c.gWaitGroup.Add(1)
		go func(i int) {
			infos := GenGallaryInfos(i)
			for _, info := range infos {
				d := GenHonDetails(info)
				if d.PageNum > 500 {
					continue
				}
				if !parseTages(d.Tags) {
					continue
				}
				if settings.CrawlerSetting.TagConsumerCount > 0 {
					for _, t := range d.Tags {
						c.sendTag(t)
					}
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
				if !c.tagSet.has(tag) {
					c.tagSet.insert(tag)
					SaveTag(tag)
				}
				c.mutex.Unlock()
			}
		}()
	}
}

func (c *coordinator) Start() {
	c.consumeHon(settings.CrawlerSetting.HonConsumerCount)
	if settings.CrawlerSetting.TagConsumerCount > 0 {
		c.comsumeTag(settings.CrawlerSetting.TagConsumerCount)
	}
	c.generateHon(
		settings.CrawlerSetting.PageStart,
		settings.CrawlerSetting.PageEnd,
	)

	c.gWaitGroup.Wait()
	close(c.honChannel)
	close(c.tagChannel)
	c.dWaitGroup.Wait()
}
