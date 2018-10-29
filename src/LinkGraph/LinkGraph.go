package LinkGraph

/*************************************
Link Graph data structures and methods
*************************************/
type LinkNode struct {
	Url			string
	ChildLinks	[]string
}

func NewLinkNode(url string) LinkNode {
	NewNode := LinkNode{
		Url: url,
		ChildLinks: []string{},
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