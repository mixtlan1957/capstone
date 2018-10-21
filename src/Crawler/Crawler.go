package Crawler

import (
	"fmt"
	"net/http"
	"strings"
	// "time"

	"LinkGraph"
	"Queue"

	"golang.org/x/net/html"

	"gopkg.in/mgo.v2"
    // "gopkg.in/mgo.v2/bson"
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

/* Structure of DB entry

	CrawlId : (bfs|dfs)Crawl_unixTimeStamp,
	RootUrl : <root_url>,
	LinkData : [
		Url : <url name>,
		Visited : <bool>,
		SqlInjectionData : <string>,
		XssData : <string>,
		ChildLinks : []<string>,
	],

*/

type CrawlDBEntry struct {
	CrawlId		string
	RootUrl		string
	LinkData	[]LinkGraph.LinkNode
}

func InsertCrawlResultsIntoDB(crawlResults map[string]*LinkGraph.LinkNode, rootUrl string) {
	session, _ := mgo.Dial("localhost:27017")
	defer session.Close()

	conn := session.DB("crawlResults").C("bfs")

	_ = conn.Insert(&CrawlDBEntry{
		CrawlId: "bfsCrawl_10/21/18_4324325346523",
		RootUrl: rootUrl,
		LinkData: []LinkGraph.LinkNode{
			LinkGraph.LinkNode{
				Url: "http://google.com",
				Visited: false,
				ChildLinks: []string{
					"http://mail.google.com",
					"http://maps.google.com",
				},
			},
			LinkGraph.LinkNode{
				Url: "http://yahoo.com",
				Visited: true,
				ChildLinks: []string{
					"http://mail.yahoo.com",
					"http://maps.yahoo.com",
				},
			},
		},
	})
}

// Breadth first search (takes root URL)
func WebCrawlBreadthFirstSearch(startUrl string) {

	// Root node for Queue
	RootUrlNode := LinkGraph.NewLinkNode(startUrl)

	// The graph of links, used to store links in DB at end
	UrlGraph := LinkGraph.CreateLinkGraph()

	// Queue of link nodes
	CrawlerQueue := Queue.NewQueue()

	// Add the root to the queue and the graph
	Queue.Enqueue(&RootUrlNode, &CrawlerQueue)
	LinkGraph.AddLinkToGraph(UrlGraph, &RootUrlNode)

	fmt.Println("Starting crawl...")

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

	fmt.Println("Crawl finished")
	for k, v := range UrlGraph {
		fmt.Printf("%v was visited : %t\nChild links:\n", k, v.Visited)

		for _, b := range v.ChildLinks {
			fmt.Printf("\t%v\n", b)
		}

		fmt.Println()
	}

	InsertCrawlResultsIntoDB(UrlGraph, startUrl)
}