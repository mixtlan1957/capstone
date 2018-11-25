package Crawler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"LinkGraph"
	"Queue"

	"golang.org/x/net/html"

	"gopkg.in/mgo.v2"
)

var wg sync.WaitGroup

type CrawlDBEntry struct {
	CrawlId				string
	LinkData			[]LinkGraph.LinkNode
	RootUrl				string
	Timestamp			int
}

type htmlForm struct {
	action	string
	inputs	[]string
	method	string
}

// Grab all of the links on a web page
func GetPageLinks(url string, baseUrl string) ([]string, []htmlForm) {
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
	parsingForm := false
	extractedForms := []htmlForm{}

	for tokType := htmlReader.Next(); tokType != html.ErrorToken; {

		tok := htmlReader.Token()
		
		// Examine the a elements (tok.DataAtom == 1) to get the link
		if tokType == html.StartTagToken && int(tok.DataAtom) == 1 {		
			// Get the href attribute value, then make sure it's a local 
			// link (ie, ones that don't begin with http/https), then 
			// append these to the list of links
			for i := 0; i < len(tok.Attr); i++ {

				if tok.Attr[i].Key == "href" && len(tok.Attr[i].Val) > 0 {

					splitLink := strings.Split(tok.Attr[i].Val, "#")
					if len(splitLink[0]) > 0 && !strings.HasPrefix(splitLink[0], "http") && !strings.HasPrefix(splitLink[0], "www.") {
						basePart := baseUrl
						if baseUrl[len(baseUrl)-1] != '/' {
							basePart += "/"
						}

						endPart := splitLink[0]
						if len(endPart) > 1 {
							for len(endPart) > 1 && (endPart[0] == '/' || endPart[0] == '.' || endPart[0] == ' ') {
								endPart = endPart[1:]
							}
						} else if len(endPart) == 1 {
							if endPart == "/" || endPart == "." || endPart == " " {
								endPart = ""
							}
						}

						if !strings.HasPrefix(endPart, "javascript:") {
							linkParts := []string{basePart, endPart}
							newLink := strings.Join(linkParts, "")
							links = append(links, newLink)
						}
					}
				}
			}
		// If the beginning of a form has been reached
		} else if int(tok.DataAtom) == 159236 {
			if (tokType == html.StartTagToken) {
				parsingForm = true
				extractedForms = append(extractedForms, htmlForm{})

				for i := 0; i < len(tok.Attr); i++ {
					switch (tok.Attr[i].Key) {
					case "action":
						extractedForms[len(extractedForms)-1].action = tok.Attr[i].Val
					case "method":
						extractedForms[len(extractedForms)-1].method = tok.Attr[i].Val
					}
				}
			} else {
				parsingForm = false
			}
		// If a form is being processed and an input is reached
		} else if parsingForm && int(tok.DataAtom) == 281349 {
			for i := 0; i < len(tok.Attr); i++ {
				switch (tok.Attr[i].Key) {
				case "name":
					extractedForms[len(extractedForms)-1].inputs = append(extractedForms[len(extractedForms)-1].inputs, tok.Attr[i].Val)
				}
			}
		}

		tokType = htmlReader.Next()
	}

	// Return the list of retrieved links
	return links, extractedForms
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

	fmt.Printf(newCrawlEntry.CrawlId)
}

func sqlInjectionFuzz(link string, forms []htmlForm) bool {
	evilPayload := "e' or 1=1; --"	// Malicious payload
	dummyPayload := "maroongolf"	// Dummy string
	isVulnerable := false			// Indicator for whether page is vulnerable to SQLi attacks

	// For each form, put sql injection in each input, make malicious request
	for form := 0; form < len(forms); form++ {

		// For each form field, fill with a malicious payload, fill the rest
		// of the form fields with the dummy values, then make the request.
		// This is to both test the form and each individual form field
		for input := 0; input < len(forms[form].inputs); input++ {
			formVals := url.Values{}

			for name := 0; name < len(forms[form].inputs); name++ {
				if input == name {
					formVals.Add(forms[form].inputs[name], evilPayload)
				} else {
					formVals.Add(forms[form].inputs[name], dummyPayload)
				}
			}

			baseUrl := link
			if link[len(link)-1] != '/' {
				baseUrl += "/"
			}
			formAction := forms[form].action
			if len(formAction) > 1 {
				for formAction[0] == '/' || formAction[0] == '.' || formAction[0] == ' '{
					formAction = formAction[1:]
				}
			} else if len(formAction) == 1 {
				if formAction == "/" || formAction == "." || formAction == " " {
					formAction = ""
				}
			}

			postUrlParts := []string{baseUrl, formAction}
			postUrl := strings.Join(postUrlParts, "")
			response, _ := http.PostForm(postUrl, formVals)
			responseBuffer := new(bytes.Buffer)
			responseBuffer.ReadFrom(response.Body)
			resStr := responseBuffer.String()

			isVulnerable = strings.Contains(resStr, "You have an error in your SQL syntax")
		}
	}

	wg.Done()

	// Return vulnerability status
	return isVulnerable
}

// Source: https://www.geeksforgeeks.org/depth-first-search-or-dfs-for-a-graph/
func DepthFirstSearch(visitedUrlMap map[string]*LinkGraph.LinkNode, node *LinkGraph.LinkNode, rootUrl string) {
	wg.Add(2)
	
	// Get child links from the parent
	nodeLinks, forms := GetPageLinks(node.Url, rootUrl)

	// SQL injection fuzz the link
	isSqliVulnerable := false
	go func() {
		isSqliVulnerable = sqlInjectionFuzz(rootUrl, forms)
	}()

	// XSS testing placeholder
	go func() {
		wg.Done()
	}()

	wg.Wait()

	// Mark link as visited
	node.SqliVulnerable = isSqliVulnerable
	LinkGraph.AddLinkToVisited(visitedUrlMap, node)

	// Add each discovered link to list of child links for parents, 
	// then if not visited continue DFS on that link
	for link := 0; link < len(nodeLinks); link++ {

		// Add child link to list of children for parent
		newNode := LinkGraph.NewLinkNode(nodeLinks[link])
		LinkGraph.AddChildLinkToParent(&newNode, node)
	
		// Continue DFS on link if not visited
		if visitedUrlMap[newNode.Url] == nil {
			DepthFirstSearch(visitedUrlMap, &newNode, rootUrl)
		}
	}	
}

// Source: https://www.geeksforgeeks.org/depth-first-search-or-dfs-for-a-graph/
func DepthFirstSearchCrawl(startUrl string) {
	// Root node
	rootUrlNode := LinkGraph.NewLinkNode(startUrl)

	// Map of visited links
	visitedUrlMap := LinkGraph.CreateLinkGraph()

	// DFS
	DepthFirstSearch(visitedUrlMap, &rootUrlNode, startUrl)

	// Store the crawl data in the DB
	InsertCrawlResultsIntoDB("dfsCrawl", visitedUrlMap, startUrl)
}

// Breadth first search (takes root URL)
// Source: https://www.geeksforgeeks.org/breadth-first-search-or-bfs-for-a-graph/
func BreadthFirstSearchCrawl(startUrl string) {
	// Root node for Queue
	rootUrlNode := LinkGraph.NewLinkNode(startUrl)

	// The graph of visited links, used to keep track of
	// which links have been visited as well as to
	// store links in DB at end
	visitedUrlMap := LinkGraph.CreateLinkGraph()

	// Queue of link nodes
	crawlerQueue := Queue.NewQueue()

	// Add the root to the queue and the graph
	Queue.Enqueue(&rootUrlNode, &crawlerQueue)
	LinkGraph.AddLinkToVisited(visitedUrlMap, &rootUrlNode)

	// While Queue isn't empty
	for crawlerQueue.Size > 0 {
		wg.Add(2)

		// Dequeue from queue
		nextQueueNode := Queue.Dequeue(&crawlerQueue)

		// Get Links from the dequeued link
		nodeLinks, forms := GetPageLinks(nextQueueNode.Url, startUrl)

		// SQL injection fuzzing
		go func() {
			visitedUrlMap[nextQueueNode.Url].SqliVulnerable = sqlInjectionFuzz(startUrl, forms)
		}()

		// XSS testing placeholder
		go func() {
			wg.Done()
		}()

		wg.Wait()

		// For each link, if not visited, enqueue link
		for link := 0; link < len(nodeLinks); link++ {

			// If link not in map of parent child links, add
			newNode := LinkGraph.NewLinkNode(nodeLinks[link])
			LinkGraph.AddChildLinkToParent(&newNode, nextQueueNode)
		
			// Enqueue unvisited link, add to map of visited links
			if visitedUrlMap[newNode.Url] == nil {
				LinkGraph.AddLinkToVisited(visitedUrlMap, &newNode)
				Queue.Enqueue(&newNode, &crawlerQueue)
			}
		}
	}

	// Store the crawl data in the DB
	InsertCrawlResultsIntoDB("bfsCrawl", visitedUrlMap, startUrl)
}