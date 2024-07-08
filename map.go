package wire

import (
	"slices"
)

// MapVal implements a Map type which might also be a multimap.
type MapVal[K comparable, V any] []Tup2Val[K, V]

// Sort the Map based on the given cmp function which may implement multimap sorting.
func (m MapVal[K, V]) Sort(cmp func(a, b Tup2Val[K, V]) int) {
	slices.SortFunc(m, cmp)
}

// Compact sorts the map entries and compacts equal keys.
//
// Compact eliminates multimap properties.
func (m *MapVal[K, V]) Compact(cmp func(a, b Tup2Val[K, V]) int) {
	m.Sort(cmp)
	*m = slices.CompactFunc(*m, func(a, b Tup2Val[K, V]) bool { return a.E0 == b.E0 })
}

func RawMap[K comparable, V any](key Proto[K], val Proto[V]) ProtoRanger[[]Tup2Val[K, V], Tup2Val[K, V]] {
	return RawSeq(Tup2(key, val))
}

func Map[K comparable, V any](key Proto[K], val Proto[V]) ProtoMakeRanger[[]Tup2Val[K, V], SpanElem[[]Tup2Val[K, V]], Tup2Val[K, V]] {
	return spanMakeRanger[[]Tup2Val[K, V], Tup2Val[K, V]](RawMap(key, val))
}
