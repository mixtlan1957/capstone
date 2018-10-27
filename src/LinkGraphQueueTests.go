package main

import . "fmt"
import . "LinkGraph"
import . "Queue"

func main() {
	/***************************************
	LinkGraph unit tests
	***************************************/
	a := LinkNode{
		Url: "http://google.com",
		Visited: false,
		ChildLinks: nil,
	}

	b := LinkNode{
		Url: "http://facebook.com",
		Visited: false,
		ChildLinks: nil,
	}

	c := LinkNode{
		Url: "http://yahoo.com",
		Visited: false,
		ChildLinks: nil,
	}

	AddChildLinkToParent(&b, &a)
	AddChildLinkToParent(&c, &a)
	AddChildLinkToParent(&b, &c)

	Printf("\n****** LinkGraph unit tests ******\n")
	Printf("Link node a has Url 'http://google.com'   : %t\n", a.Url == "http://google.com")
	Printf("Link node b has Url 'http://facebook.com' : %t\n", b.Url == "http://facebook.com")
	Printf("Link node c has Url 'http://yahoo.com'    : %t\n", c.Url == "http://yahoo.com")
	Println()

	Printf("Link node a has child link b : %t\n", a.ChildLinks[0] == &b)
	Printf("Link node a has child link c : %t\n", a.ChildLinks[1] == &c)
	Printf("Link node c has child link b : %t\n", c.ChildLinks[0] == &b)
	Println()

	graph := CreateLinkGraph()

	AddLinkToGraph(graph, &a)
	AddLinkToGraph(graph, &b)
	AddLinkToGraph(graph, &c)
	Printf("Link node for 'http://google.com'   is a : %t\n", graph["http://google.com"] == &a)
	Printf("Link node for 'http://facebook.com' is b : %t\n", graph["http://facebook.com"] == &b)
	Printf("Link node for 'http://yahoo.com'    is c : %t\n", graph["http://yahoo.com"] == &c)
	Println()

	/***************************************
	Queue unit tests
	***************************************/
	q := NewQueue()

	Printf("\n****** Queue unit tests ******\n")

	Printf("Queue Back is nil                     : %t\n", q.Back == nil)
	Enqueue(&a, &q)
	Printf("Queue Back is a (http://google.com)   : %t\n", q.Back.Val.Url == "http://google.com")
	Enqueue(&b, &q)
	Printf("Queue Back is b (http://facebook.com) : %t\n", q.Back.Val.Url == "http://facebook.com")
	Enqueue(&c, &q)
	Printf("Queue Back is c (http://yahoo.com)    : %t\n", q.Back.Val.Url == "http://yahoo.com")
	Println()

	bn := Dequeue(&q)
	Printf("Queue Front was a (http://google.com) : %t\n", bn.Url == "http://google.com")
	bn = Dequeue(&q)
	Printf("Queue Front was b (http://google.com) : %t\n", bn.Url == "http://facebook.com")
	bn = Dequeue(&q)
	Printf("Queue Front was c (http://google.com) : %t\n", bn.Url == "http://yahoo.com")
	bn = Dequeue(&q)
	Printf("Queue Front is nil                    : %t\n", bn == nil)
	Println()
}