package main

import (
	"Crawler"
)

func main() {
	//url := "http://192.168.1.158/mutillidae/"
	url := "http://localhost:8000/"
	
	Crawler.BreadthFirstSearchCrawl(url)
	Crawler.DepthFirstSearchCrawl(url)
}


