package runner

import "cmp"

type flatKV[K, V cmp.Ordered] struct {
	k K
	v V
}

// flattenedKVs can be used to convert a map into a sorted slice of key-value pairs
type flattenedKVs[K, V cmp.Ordered] []flatKV[K, V]

func (flattenedKVs[K, V]) cmpV(desc bool) func(a, b flatKV[K, V]) int {
	return func(a, b flatKV[K, V]) int {
		if desc {
			a, b = b, a
		}
		return cmp.Compare(a.v, b.v)
	}
}

func (flattenedKVs[K, V]) cmpK(desc bool) func(a, b flatKV[K, V]) int {
	return func(a, b flatKV[K, V]) int {
		if desc {
			a, b = b, a
		}
		return cmp.Compare(a.k, b.k)
	}
}

func flatten[K, V cmp.Ordered](m map[K]V) flattenedKVs[K, V] {
	out := make([]flatKV[K, V], 0, len(m))
	for k, v := range m {
		out = append(out, flatKV[K, V]{k, v})
	}
	return out
}
