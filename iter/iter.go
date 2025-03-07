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

func KeyMap[K comparable](iter iter.Seq[K]) iter.Seq2[K, struct{}] {
	return func(yield func(K, struct{}) bool) {
		for k := range iter {
			if !yield(k, struct{}{}) {
				return
			}
		}
	}
}

func ValueMap[V any](iter iter.Seq[V]) iter.Seq2[struct{}, V] {
	return func(yield func(struct{}, V) bool) {
		for v := range iter {
			if !yield(struct{}{}, v) {
				return
			}
		}
	}
}
