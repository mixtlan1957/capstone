package Crawler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"LinkGraph"
	"Queue"

	"golang.org/x/net/html"

	"gopkg.in/mgo.v2"
)

// Grab all of the links on a web page
func GetPageLinks(url string, baseUrl string) []string {
	// Get request to the URL
	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	defer response.Body.Close()

	var links []string
	// src: https://mionskowski.pl/html-golang-stream-processing
	// Iterate through each element of the HTML response
	htmlReader := html.NewTokenizer(response.Body)
	for tokType := htmlReader.Next(); tokType != html.ErrorToken; {

		tok := htmlReader.Token()

		// Examine the a elements (tok.DataAtom == 1)
		if tokType == html.StartTagToken && int(tok.DataAtom) == 1 {		
			// Get the href attribute value, then make sure it's a local 
			// link (ie, ones that don't begin with http/https), then 
			// append these to the list of links
			for i := 0; i < len(tok.Attr); i++ {

				if tok.Attr[i].Key == "href" && len(tok.Attr[i].Val) > 0 {

					splitLink := strings.Split(tok.Attr[i].Val, "#")
					if len(splitLink[0]) > 0 && !strings.HasPrefix(splitLink[0], "http") {

						linkParts := []string{baseUrl, splitLink[0]}
						newLink := strings.Join(linkParts, "")
						links = append(links, newLink)
					}
				}
			}
		}

		tokType = htmlReader.Next()
	}

	// Return the list of retrieved links
	return links
}

type CrawlDBEntry struct {
	CrawlId				string
	LinkData			[]LinkGraph.LinkNode
	RootUrl				string
	Timestamp			int
}

// Takes the crawlResults from the crawl and inserts into the appropriate "crawlCollection" collection in the db
// SOURCE: https://labix.org/mgo
func InsertCrawlResultsIntoDB(crawlCollection string, crawlResults map[string]*LinkGraph.LinkNode, rootUrl string) {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	conn := session.DB("crawlResults").C(crawlCollection)

	// Unique timestamp for the crawl
	crawlTimestamp := int(time.Now().Unix())
	
	// Create a crawl entry struct
	newCrawlEntry := CrawlDBEntry{
		CrawlId: strings.Join([]string{crawlCollection, rootUrl, strconv.Itoa(crawlTimestamp)}, "_"),
		LinkData: nil,
		RootUrl: rootUrl,
		Timestamp: crawlTimestamp,
	}

	// Append all the found links and their data to the crawl entry struct LinkData
	for _, node := range crawlResults {
		newCrawlEntry.LinkData = append(newCrawlEntry.LinkData, *node)
	}

	// Insert the crawl into the db
	err = conn.Insert(&newCrawlEntry)
	if err != nil {
		fmt.Println(err)
	}
}

func DepthFirstSearch(urlMap map[string]*LinkGraph.LinkNode, node *LinkGraph.LinkNode, rootUrl string) {
	// Get Links from the dequeued link
	nodeLinks := GetPageLinks(node.Url, rootUrl)

	// Mark Link as visited
	node.Visited = true

	// For each link, if not visited, then visit
	for link := 0; link < len(nodeLinks); link++ {

		// If link not in map of parent child links, add
		newNode := LinkGraph.NewLinkNode(nodeLinks[link])
		LinkGraph.AddChildLinkToParent(&newNode, node)
	
		// Add to graph if not already in and do DFS on link
		if urlMap[newNode.Url] == nil {
			LinkGraph.AddLinkToGraph(urlMap, &newNode)
			
			if newNode.Visited == false {
				DepthFirstSearch(urlMap, &newNode, rootUrl)
			}
		}
	}	
}

func DepthFirstSearchCrawl(startUrl string) {
	// Root node
	RootUrlNode := LinkGraph.NewLinkNode(startUrl)

	// The graph of links, used to store links in DB at end
	UrlGraph := LinkGraph.CreateLinkGraph()

	// DFS
	DepthFirstSearch(UrlGraph, &RootUrlNode, startUrl)

	// Display the crawl results to the console
	fmt.Println("Crawl finished\n")
	for k, v := range UrlGraph {
		fmt.Printf("%v was visited : %t\nChild links:\n", k, v.Visited)

		for _, b := range v.ChildLinks {
			fmt.Printf(" * %v\n", b)
		}

		fmt.Println()
	}

	// Store the crawl data in the DB
	InsertCrawlResultsIntoDB("bfsCrawl", UrlGraph, startUrl)
}

// Breadth first search (takes root URL)
func BreadthFirstSearchCrawl(startUrl string) {

	// Root node for Queue
	RootUrlNode := LinkGraph.NewLinkNode(startUrl)

	// The graph of links, used to store links in DB at end
	UrlGraph := LinkGraph.CreateLinkGraph()

	// Queue of link nodes
	CrawlerQueue := Queue.NewQueue()

	// Add the root to the queue and the graph
	Queue.Enqueue(&RootUrlNode, &CrawlerQueue)
	LinkGraph.AddLinkToGraph(UrlGraph, &RootUrlNode)

	fmt.Println("Starting crawl...\n")

	// While Queue isn't empty
	for CrawlerQueue.Size > 0 {

		// Dequeue from queue
		nextQueueNode := Queue.Dequeue(&CrawlerQueue)
		
		if !nextQueueNode.Visited {

			// Security scan
			//fmt.Printf("Scanning %v\n", nextQueueNode.Url)

			// Get Links from the dequeued link
			nodeLinks := GetPageLinks(nextQueueNode.Url, startUrl)

			// Mark Link as visited
			nextQueueNode.Visited = true

			// For each link, if not visited, enqueue link
			for link := 0; link < len(nodeLinks); link++ {

				// If link not in map of parent child links, add
				newNode := LinkGraph.NewLinkNode(nodeLinks[link])
				LinkGraph.AddChildLinkToParent(&newNode, nextQueueNode)
			
				// Add to graph if not already in and enqueue
				if UrlGraph[newNode.Url] == nil {
					LinkGraph.AddLinkToGraph(UrlGraph, &newNode)
					Queue.Enqueue(&newNode, &CrawlerQueue)
				}
			}
		}
	}

	// Display the crawl results to the console
	fmt.Println("Crawl finished\n")
	for k, v := range UrlGraph {
		fmt.Printf("%v was visited : %t\nChild links:\n", k, v.Visited)

		for _, b := range v.ChildLinks {
			fmt.Printf(" * %v\n", b)
		}

		fmt.Println()
	}

	// Store the crawl data in the DB
	InsertCrawlResultsIntoDB("bfsCrawl", UrlGraph, startUrl)
}