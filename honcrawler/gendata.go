package honcrawler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

/*
	This `go` file will generate the details of ero-hons and collect
	there tags.
*/

// from https://www.wnacg.com/albums-index-page-1.html
// to https://www.wnacg.com/albums-index-page-7101.html

const Host = `https://www.wnacg.com`
const GallaryUrl string = Host + `/albums-index-page-%d.html`
const UserAgent string = `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:104.0) Gecko/20100101 Firefox/104.0`
const ImgsPerPage int = 12

// page number regex
var pattern = regexp.MustCompile(`(\d+)P`)

// typical HonUrl /photos-index-aid-169728.html
type GallaryInfo struct {
	HonUrl string
	Title  string
}

type HonDetail struct {
	Tags    []string
	Title   string
	PageNum int      // 可以计算翻页次数
	Images  []string // 按顺序的本子页面Url
}

func GenGallaryInfos(page int) (infos []*GallaryInfo) {
	GallaryCollector := collector.Clone()
	GallaryCollector.OnHTML(".pic_box>a", func(e *colly.HTMLElement) {
		info := &GallaryInfo{}
		info.HonUrl = e.Attr("href")
		info.Title = e.Attr("title")
		infos = append(infos, info)
	})

	GallaryCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Requesting:", r.URL)
	})

	url := fmt.Sprintf(GallaryUrl, page)
	GallaryCollector.Visit(url)
	return
}

func GenHonDetails(g *GallaryInfo) (details *HonDetail) {
	details = &HonDetail{
		Tags:   make([]string, 0),
		Title:  g.Title,
		Images: make([]string, 0),
	}

	details.crawlTagAndPage(g)
	details.crawlImages(g)

	return
}

func (hd *HonDetail) crawlTagAndPage(g *GallaryInfo) {
	HonCollector := collector.Clone()
	// 对于标题信息的处理，获取Tags和PageNum
	HonCollector.OnHTML(".uwconn", func(e *colly.HTMLElement) {
		e.ForEach(".uwconn>label", func(i int, h *colly.HTMLElement) {
			if i == 0 { // 分类，解析到tag里，从0开始
				tags := strings.Split(h.Text, " / ")
				hd.Tags = append(hd.Tags, tags...)
			}
			if i == 1 { // 页数，解析到PageNum
				pageStr := pattern.FindAllStringSubmatch(h.Text, -1)
				cnt, err := strconv.Atoi(pageStr[0][1])
				if err != nil {
					fmt.Printf("wrong with str unmarshal %v", pageStr)
				}
				hd.PageNum = cnt
			}
		})
		e.ForEach(".tagshow", func(_ int, h *colly.HTMLElement) {
			hd.Tags = append(hd.Tags, h.Text)
		})
	})
	HonCollector.Visit(Host + g.HonUrl)
}

func (hd *HonDetail) crawlImages(g *GallaryInfo) {
	HonCollector := collector.Clone()
	total := hd.PageNum/ImgsPerPage + 1
	HonCollector.OnHTML(".pic_box>a", func(e *colly.HTMLElement) {
		// e.ForEach("a", func(i int, h *colly.HTMLElement) {

		// })
		hd.Images = append(hd.Images, e.Attr("href"))
	})

	for i := 1; i <= total; i++ {
		url := pageUrlTrans(g.HonUrl, i)
		HonCollector.Visit(Host + url)
	}
}

// `/photos-index-aid-169728.html`` =>`/photos-index-page-i-aid-169728.html`
func pageUrlTrans(u string, i int) string {
	strs := strings.Split(u, "-")
	ret := fmt.Sprintf("%s-%s-page-%d-%s-%s", strs[0], strs[1], i, strs[2], strs[3])
	return ret
}
