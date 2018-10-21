package Queue

import . "LinkGraph"

/*************************************
Queue data structures and methods
*************************************/
type QueueNode struct {
	Val			*LinkNode
	Next		*QueueNode
}

type Queue struct {
	Front		*QueueNode
	Back		*QueueNode
	Size		int
}

func NewQueue() Queue {
	return Queue{
		Front: nil,
		Back: nil,
		Size: 0,
	}
}

func Enqueue(node *LinkNode, q *Queue) {
	newQueueNode := QueueNode{
		Val: node,
		Next: nil,
	}
	
	if q.Size == 0 {
		q.Front = &newQueueNode
		q.Back = &newQueueNode
	} else {
		q.Back.Next = &newQueueNode
		q.Back = &newQueueNode
	}

	q.Size++
}

func Dequeue(q *Queue) *LinkNode {
	if q.Size == 0 {
		return nil
	}

	frontNode := q.Front
	q.Front = q.Front.Next

	q.Size--

	return frontNode.Val
}