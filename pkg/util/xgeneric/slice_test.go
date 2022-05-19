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
