package main

import (
	"fmt"

	"github.com/Youngkingman/gentlemanSpider/honcrawler"
	"github.com/Youngkingman/gentlemanSpider/settings"
)

func main() {
	fmt.Println(*settings.CrawlerSetting)
	honcrawler.Coordinator.Start()
}
