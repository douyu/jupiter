package xgeneric

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
func Chunk[T comparable](collection []T, size int) [][]T {
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
