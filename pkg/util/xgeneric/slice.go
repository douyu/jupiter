package xgeneric

import "math/rand"

// Uniq returns a duplicate-free version of an array
func Uniq[T comparable](collection []T) []T {
	result := make([]T, 0, len(collection))
	seen := make(map[T]struct{}, len(collection))

	for _, item := range collection {
		if _, ok := seen[item]; ok {
			continue
		}

		seen[item] = struct{}{}
		result = append(result, item)
	}

	return result
}

// Chunk returns an array of elements split into groups the length of size.
func Chunk[T any](collection []T, size int) [][]T {
	if size <= 0 {
		panic("size must bigger than 0")
	}

	result := make([][]T, 0, len(collection)/2+1)
	length := len(collection)

	for i := 0; i < length; i++ {
		chunk := i / size

		if i%size == 0 {
			result = append(result, make([]T, 0, size))
		}

		result[chunk] = append(result[chunk], collection[i])
	}

	return result
}

// Map manipulates a slice and transforms it to a slice of another type.
func Map[T any, R any](collection []T, iteratee func(T, int) R) []R {
	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = iteratee(item, i)
	}

	return result
}

// Subset return part of a slice.
func Subset[T any](collection []T, offset int, limit uint) []T {
	size := len(collection)

	if offset < 0 {
		offset = size + offset
		if offset < 0 {
			offset = 0
		}
	}

	if offset > size {
		return []T{}
	}

	if limit > uint(size)-uint(offset) {
		limit = uint(size - offset)
	}

	return collection[offset : offset+int(limit)]
}

// Shuffle returns an array of shuffled values. Uses the Fisher-Yates shuffle algorithm.
func Shuffle[T any](collection []T) []T {
	rand.Shuffle(len(collection), func(i, j int) {
		collection[i], collection[j] = collection[j], collection[i]
	})

	return collection
}

// Filter iterates over elements of collection, returning an array of all elements predicate returns truthy for.
func Filter[V any](collection []V, predicate func(V, int) bool) []V {
	result := []V{}

	for i, item := range collection {
		if predicate(item, i) {
			result = append(result, item)
		}
	}

	return result
}
