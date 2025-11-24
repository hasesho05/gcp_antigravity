package util

import "github.com/samber/lo"

// Map is a helper function for lo.Map when the index is not needed.
func Map[T any, R any](collection []T, iteratee func(item T) R) []R {
	return lo.Map(collection, func(item T, index int) R {
		return iteratee(item)
	})
}
