package xgeneric

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniq(t *testing.T) {
	result := Uniq([]int{1, 2, 2, 3, 3, 4, 4})

	assert.Equal(t, len(result), 4)
	assert.Equal(t, result, []int{1, 2, 3, 4})

	result1 := Uniq([]int64{1, 2, 2, 3, 3, 4, 4})

	assert.Equal(t, len(result1), 4)
	assert.Equal(t, result1, []int64{1, 2, 3, 4})
}

func TestChunk(t *testing.T) {
	result := Chunk([]int{1, 2, 2, 3, 3, 4, 4}, 2)

	assert.Equal(t, len(result), 4)
	assert.Equal(t, result, [][]int([][]int{{1, 2}, {2, 3}, {3, 4}, {4}}))

	assert.PanicsWithValue(t, "size must bigger than 0", func() {
		Chunk([]int{1, 2, 3}, 0)
	})
}

func TestMap(t *testing.T) {
	result := Map([]int{1, 2, 2, 3, 3, 4, 4}, func(t int, i int) string {
		return strconv.Itoa(t)
	})

	assert.Equal(t, len(result), 7)
	assert.Equal(t, result, []string([]string{"1", "2", "2", "3", "3", "4", "4"}))
}

func TestSubset(t *testing.T) {
	result := Subset([]int{1, 2, 2, 3, 3, 4, 4}, 0, 4)

	assert.Equal(t, len(result), 4)
	assert.Equal(t, result, []int([]int{1, 2, 2, 3}))

	assert.Equal(t, Subset([]int{1, 2, 2, 3, 3, 4, 4}, 4, 4), []int([]int{3, 4, 4}))
	assert.Equal(t, Subset([]int{1, 2, 2, 3, 3, 4, 4}, 8, 4), []int([]int{}))

	assert.Equal(t, Subset([]int{1, 2, 2, 3, 3, 4, 4}, -1, 4), []int([]int{4}))
	assert.Equal(t, Subset([]int{1, 2, 2, 3, 3, 4, 4}, -9, 4), []int{1, 2, 2, 3})
}

func TestShuffle(t *testing.T) {

	result1 := Shuffle([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	result2 := Shuffle([]int{})

	assert.NotEqual(t, result1, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	assert.Equal(t, result2, []int{})
}

func TestFilter(t *testing.T) {

	result1 := Filter([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, func(v int, i int) bool {
		return v == 8
	})
	result2 := Filter([]int{}, func(v int, i int) bool {
		return false
	})

	assert.NotEqual(t, result1, []int{0, 1, 2, 3, 4, 5, 6, 7, 9, 10})
	assert.Equal(t, result2, []int{})
}
