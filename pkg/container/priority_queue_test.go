package container

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue()
	elements := []int{5, 3, 7, 8, 6, 2, 9}
	for _, e := range elements {
		assert.Nil(t, pq.Push(e, e))
	}

	sort.Ints(elements)
	for _, e := range elements {
		item, err := pq.Pop()
		assert.Nil(t, err)
		assert.Equal(t, e, item.(int))
	}
}

func TestPriorityQueueUpdate(t *testing.T) {
	pq := NewPriorityQueue()
	assert.Nil(t, pq.Push("foo", 3))
	assert.Nil(t, pq.Push("bar", 4))
	pq.UpdatePriority("bar", 2)

	item, err := pq.Pop()
	assert.Nil(t, err)
	assert.Equal(t, item.(string), "bar")
}

func TestPriorityQueueLen(t *testing.T) {
	pq := NewPriorityQueue()
	assert.Equal(t, 0, pq.Len())
	assert.Nil(t, pq.Push("foo", 1))
	assert.Nil(t, pq.Push("bar", 1))
	assert.Equal(t, 2, pq.Len())
}

func TestDoubleAddition(t *testing.T) {
	pq := NewPriorityQueue()
	assert.Nil(t, pq.Push("foo", 2))
	assert.Nil(t, pq.Push("bar", 3))
	assert.Nil(t, pq.Push("bar", 1))
	assert.Equal(t, 2, pq.Len())

	item, _ := pq.Pop()
	assert.Equal(t, item.(string), "foo")
}

func TestPopEmptyQueue(t *testing.T) {
	pq := NewPriorityQueue()
	_, err := pq.Pop()
	assert.NotNil(t, err)
}

func TestUpdateNonExistingItem(t *testing.T) {
	pq := NewPriorityQueue()

	assert.Nil(t, pq.Push("foo", 4))
	pq.UpdatePriority("bar", 5)
	assert.Equal(t, 1, pq.Len())

	item, _ := pq.Pop()
	assert.Equal(t, item.(string), "foo")
}
