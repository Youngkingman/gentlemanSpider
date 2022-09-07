package honcrawler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// 从 https://www.wnacg.com/albums-index-page-1.html
// 爬到 https://www.wnacg.com/albums-index-page-7094.html
const Host = `https://www.wnacg.com`
const GallaryUrl string = Host + `/albums-index-page-%d.html`
const UserAgent string = `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:104.0) Gecko/20100101 Firefox/104.0`
const ImgsPerPage int = 12

// 用于收集画廊中的本子链接与标题信息
var GallaryCollector = colly.NewCollector(
	colly.UserAgent(UserAgent),
	colly.AllowURLRevisit(),
)

// 用于收集进入Url后的相关信息
var HonCollector = GallaryCollector.Clone()

func init() {
	GallaryCollector.SetRequestTimeout(120 * time.Second)
	// 注册一些差错处理
	GallaryCollector.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error %s: %v\n", r.Request.URL, err)
	})
}

// 一个典型的HonUrl /photos-index-aid-169728.html
type GallaryInfo struct {
	HonUrl string
	Title  string
}

func GetGallaryInfos(page int) (infos []*GallaryInfo) {
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

type HonDetail struct {
	Tags    []string
	Title   string   // 复用上层+Tas
	PageNum int      // 可以计算翻页次数
	Images  []string // 按顺序的本子页面Url
}

func GetHonDetail(g *GallaryInfo) (details *HonDetail) {
	// 对于标题信息的处理，获取Tags和PageNum
	HonCollector.OnHTML(".uwconn", func(e *colly.HTMLElement) {
		e.ForEach("label", func(i int, h *colly.HTMLElement) {
			if i == 1 { // 分类，解析到tag里
				tags := strings.Split(h.Text, "/")
				details.Tags = append(details.Tags, tags...)
			}
			if i == 2 { // 页数，解析到PageNum
				pageStr := strings.Split(h.Text, " ")[1]
				pageStr = pageStr[:len(pageStr)-1] // 去掉P
				cnt, err := strconv.Atoi(pageStr)
				if err != nil {
					fmt.Printf("wrong with str unmarshal %v", pageStr)
				}
				details.PageNum = cnt
			}
		})
		e.ForEach(".tagshow", func(_ int, h *colly.HTMLElement) {
			details.Tags = append(details.Tags, h.Text)
		})
	})
	HonCollector.Visit(Host + g.HonUrl)

	total := details.PageNum/ImgsPerPage + 1
	HonCollector.OnHTML("ul", func(e *colly.HTMLElement) {
		e.ForEach("a", func(i int, h *colly.HTMLElement) {
			details.Images = append(details.Images, h.Attr("href"))
		})
	})

	for i := 1; i <= total; i++ {
		url := pageUrlTrans(g.HonUrl, i)
		HonCollector.Visit(url)
	}
	return
}

// `/photos-index-aid-169728.html`` =>`/photos-index-page-i-aid-169728.html`
func pageUrlTrans(u string, i int) string {
	strs := strings.Split(u, "-")
	ret := fmt.Sprintf("%s-%s-page-%d-%s-%s", strs[0], strs[1], i, strs[2], strs[3])
	return ret
}
