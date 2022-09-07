package main

import (
	"fmt"

	"github.com/Youngkingman/gentlemanSpider/honcrawler"
)

func main() {
	ret := honcrawler.GetGallaryInfos(2)
	for _, v := range ret {
		fmt.Print(*v)
	}
	// ret := honcrawler.GetGallaryInfos(1)
	// t := honcrawler.GetHonDetail(ret[0])
	// fmt.Println(*t)
}
