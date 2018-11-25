package main

import (
    "Crawler"
    "flag"
    "strconv"
)

func main() {
    // Get Command Line Arguments
    url := flag.String("url", "http://localhost:8000", "The URL to crawl")
    searchType := flag.String("search", "bfs", "'bfs' for breadth first crawl, 'dfs' for depth first crawl")
    // keyword := flag.String("keyword", "test", "Specify a keyword")
    depth := flag.String("depth", "1", "Enter a depth limit 1-3")

    flag.Parse()
    depthLimit, _ := strconv.Atoi(*depth)

    if (*searchType == "bfs") {
        Crawler.BreadthFirstSearchCrawl(*url, depthLimit)
    } else {
        Crawler.DepthFirstSearchCrawl(*url, depthLimit)
    }
}