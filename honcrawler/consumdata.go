package honcrawler

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	colly "github.com/gocolly/colly/v2"
)

var patternWinFile = regexp.MustCompile(`[/\:*?"<>|]`)

func Download(hd *HonDetail) {
	UrlCollector := collector.Clone()
	ImgCollector := collector.Clone()

	/*
		While the website is continuous updating, duplicated
		hon-details may be mixin.
	*/
	outputDirTitle := "./hon/" + genDirNameAndFilter(hd) + "/"
	err := os.MkdirAll(outputDirTitle, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	ImgCollector.OnResponse(func(r *colly.Response) {
		r.Save(outputDirTitle + r.FileName())
	})

	UrlCollector.OnHTML("#picarea", func(e *colly.HTMLElement) {
		originSource := e.Attr("src")
		ImgCollector.Visit("https:" + originSource)
	})

	for _, url := range hd.Images {
		UrlCollector.Visit(Host + url)
	}
}

func SaveTag(tag string) {
	file, err := os.OpenFile("./activeTags", os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("fail crate")
	}
	write := bufio.NewWriter(file)
	write.WriteString(tag + "\r\n")
	write.Flush()
	defer file.Close()
}
