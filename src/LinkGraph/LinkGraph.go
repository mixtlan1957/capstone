package LinkGraph

/*************************************
Link Graph data structures and methods
*************************************/
type LinkNode struct {
	Url				string
	ChildLinks		[]string
	Depth			int
	HasKeyword		bool
	IsCrawlRoot		bool
	SqliVulnerable	bool
	Title			string
	XssVulnerable	bool
	TestInfo		[]string
	XssTestInfo		[]string
}

func NewLinkNode(url string) LinkNode {
	NewNode := LinkNode{
		Url: url,
		ChildLinks: []string{},
		Depth: 0,
		HasKeyword: false,
		IsCrawlRoot: false,
		SqliVulnerable: false,
		Title: "",
		XssVulnerable: false,
		TestInfo: []string{},
		XssTestInfo: []string{},
	}

	return NewNode
}

func AddChildLinkToParent(child *LinkNode, parent *LinkNode) {
	parent.ChildLinks = append(parent.ChildLinks, child.Url)
}

func CreateLinkGraph() map[string]*LinkNode {
	graph := make(map[string]*LinkNode)

	return graph
}

func AddLinkToVisited(graph map[string]*LinkNode, link *LinkNode) {
	graph[link.Url] = link
}