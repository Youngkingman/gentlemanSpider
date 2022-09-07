package main

import "github.com/Youngkingman/gentlemanSpider/honcrawler"

func main() {
	honcrawler.Coordinator.Try()
	// err := os.MkdirAll("./fuck dd/", os.ModePerm)
	// fmt.Println(err)
}
