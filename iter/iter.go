package iter

import (
	"iter"
)

func Chain[T any](iters ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, iter := range iters {
			for v := range iter {
				if !yield(v) {
					return
				}
			}
		}
	}
}
