package vibrant

import (
	"testing"
)

var priorityFunction = func(data interface{}) uint32 {
	return data.(uint32)
}

func TestDefaultPriorityQueue(t *testing.T) {
	q := NewPriorityQueue(3, priorityFunction)
	testData := [][]uint32{
		{1, 2, 3},
		{1, 3, 2},
		{2, 1, 3},
		{2, 3, 1},
		{3, 1, 2},
		{3, 2, 1},
	}
	for _, data := range testData {
		if q.Len() != 0 {
			t.Fatalf("Expected queue to be empty.\n")
		}
		for _, item := range data {
			q.Offer(item)
		}
		for expected := q.Len(); expected > 0; expected-- {
			actual := priorityFunction(q.Poll())
			if uint32(expected) != actual {
				t.Errorf("Expected item with priority %v instead of %v.\n", expected, actual)
			}
		}
	}
}

func TestDefaultPriorityQueue_MultipleOffer(t *testing.T) {
	q := NewPriorityQueue(5, priorityFunction)
	q.Offer(uint32(6), uint32(5), uint32(4), uint32(3), uint32(2), uint32(1))
	for expected := q.Len(); expected > 0; expected-- {
		actual := priorityFunction(q.Poll())
		if uint32(expected) != actual {
			t.Errorf("Expected item with priority %v instead of %v.\n", expected, actual)
		}
	}
	q.Offer(uint32(1), uint32(2), uint32(3), uint32(4), uint32(5), uint32(6))
	for expected := q.Len(); expected > 0; expected-- {
		actual := priorityFunction(q.Poll())
		if uint32(expected) != actual {
			t.Errorf("Expected item with priority %v instead of %v.\n", expected, actual)
		}
	}
	q.Offer(uint32(1), uint32(5), uint32(3), uint32(4), uint32(2), uint32(6))
	for expected := q.Len(); expected > 0; expected-- {
		actual := priorityFunction(q.Poll())
		if uint32(expected) != actual {
			t.Errorf("Expected item with priority %v instead of %v.\n", expected, actual)
		}
	}
}
