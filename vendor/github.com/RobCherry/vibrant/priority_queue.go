package vibrant

import (
	"container/heap"
)

// A PriorityQueue implements heap.Interface and holds items.
type PriorityQueue interface {
	Offer(items ...interface{})
	Poll() interface{}
	Len() int
}

// NewPriorityQueue creates a new PriorityQueue with a given capacity and priority function.
func NewPriorityQueue(initialCapacity uint32, priorityFunction func(interface{}) uint32) PriorityQueue {
	return &priorityQueue{
		make([]interface{}, 0, initialCapacity),
		priorityFunction,
	}
}

type priorityQueue struct {
	queue            []interface{}
	priorityFunction func(interface{}) uint32
}

func (pq *priorityQueue) Offer(items ...interface{}) {
	for _, item := range items {
		heap.Push(pq, item)
	}
}

func (pq *priorityQueue) Poll() interface{} {
	return heap.Pop(pq)
}

// Satisfy the heap.Interface interface.
func (pq priorityQueue) Len() int {
	return len(pq.queue)
}

// Satisfy the heap.Interface interface.
func (pq priorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq.priorityFunction(pq.queue[i]) > pq.priorityFunction(pq.queue[j])
}

// Satisfy the heap.Interface interface.
func (pq priorityQueue) Swap(i, j int) {
	pq.queue[i], pq.queue[j] = pq.queue[j], pq.queue[i]
}

// Satisfy the heap.Interface interface.
func (pq *priorityQueue) Push(item interface{}) {
	pq.queue = append(pq.queue, item)
}

// Satisfy the heap.Interface interface.
func (pq *priorityQueue) Pop() interface{} {
	old := pq.queue
	n := len(old)
	item := old[n-1]
	pq.queue = old[0 : n-1]
	return item
}
