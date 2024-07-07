package wire

import "fmt"

func RawMap[K comparable, V any](key Proto[K], val Proto[V]) Proto[[]Tup2Val[K, V]] {
	return Seq(Tup2(key, val))
}

func Map[K comparable, V any](key Proto[K], val Proto[V]) Proto[SpanElem[[]Tup2Val[K, V]]] {
	return Span(RawMap(key, val))
}

func VisitMap[K comparable, V any](key Proto[K], val Proto[V]) func(Reader, func(Tup2Val[K, V]) error) error {
	return VisitSeq(Tup2(key, val))
}

func MakeRawMap[K comparable, V any](m map[K]V) []Tup2Val[K, V] {
	res := make([]Tup2Val[K, V], 0, len(m))
	for k, v := range m {
		res = append(res, Tup2Val[K, V]{k, v})
	}
	return res
}

func MakeDeterministicRawMap[K comparable, V any](m map[K]V, keys []K) []Tup2Val[K, V] {
	res := make([]Tup2Val[K, V], 0, len(m))
	for _, k := range keys {
		v, ok := m[k]
		if !ok {
			panic(fmt.Errorf("key not found in map: %v", k))
		}
		res = append(res, Tup2Val[K, V]{k, v})
	}
	return res
}

func MakeMap[K comparable, V any](key Proto[K], val Proto[V]) func(map[K]V) SpanElem[[]Tup2Val[K, V]] {
	return func(m map[K]V) SpanElem[[]Tup2Val[K, V]] { return MakeSpan(RawMap(key, val))(MakeRawMap(m)) }
}

func MakeDeterministicMap[K comparable, V any](key Proto[K], val Proto[V]) func(map[K]V, []K) SpanElem[[]Tup2Val[K, V]] {
	return func(m map[K]V, keys []K) SpanElem[[]Tup2Val[K, V]] {
		return MakeSpan(RawMap(key, val))(MakeDeterministicRawMap(m, keys))
	}
}
