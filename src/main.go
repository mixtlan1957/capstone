package main

import (
    "Crawler"
    "flag"
)

func main() {
    // Get Command Line Arguments
    url := flag.String("url", "http://localhost:8000", "The URL to crawl")
    searchType := flag.String("search", "bfs", "'bfs' for breadth first crawl, 'dfs' for depth first crawl")
    keyword := flag.String("keyword", "", "Specify a keyword")
    depthLimit := flag.Int("depth", 2147483647, "Enter a depth limit 1-3")
    vulnerabilityScan := flag.Bool("fuzz", false, "Enter this flag")

    flag.Parse()

    if (*searchType == "bfs") {
        Crawler.BreadthFirstSearchCrawl(*url, *depthLimit, *keyword, *vulnerabilityScan)
    } else {
        Crawler.DepthFirstSearchCrawl(*url, *depthLimit, *keyword, *vulnerabilityScan)
    }
}