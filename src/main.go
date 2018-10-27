package main

import (
	"Crawler"
	"flag"
)

func main() {
	// Get Command Line Arguments
	url := flag.String("url", "http://localhost:8000", "The URL to crawl")
	searchType := flag.String("search", "bfs", "'bfs' for breadth first crawl, 'dfs' for depth first crawl")
	
	flag.Parse()

	if (*searchType == "bfs") {
		Crawler.BreadthFirstSearchCrawl(*url)
	} else {
		Crawler.DepthFirstSearchCrawl(*url)
	}
}


