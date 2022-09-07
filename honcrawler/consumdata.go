package honcrawler

import (
	"fmt"
	"os"

	colly "github.com/gocolly/colly/v2"
)

func Download(hd *HonDetail) {
	UrlCollector := collector.Clone()
	ImgCollector := collector.Clone()

	/*
		While the website is continuous updating, duplicated
		hon-details may be mixin.
	*/
	outputDirTitle := "./hon/" + genDirName(hd) + "/"
	// fmt.Printf("create output Dir %s", outputDirTitle)
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

func genDirName(hd *HonDetail) (s string) {
	s = hd.Title + "["
	for _, v := range hd.Tags {
		s = s + "_" + v
	}
	s = s + "]"
	return
}

func SaveTag(tag string) {
	//TODO wirte the tags to a file
}
