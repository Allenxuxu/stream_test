package stream

import (
	"reflect"

	"stream_test/stream/types"
)

const unknownSize = -1

type iterator interface {
	GetSizeIfKnown() int64
	HasNext() bool
	Next() types.T
}

func it(elements ...types.T) iterator {
	return &sliceIterator{
		base: &base{
			current: 0,
			size:    len(elements),
		},
		elements: elements,
	}
}

func withSupplier(get types.Supplier) iterator {
	return &supplierIt{get: get}
}

type base struct {
	current int
	size    int
}

func (b *base) GetSizeIfKnown() int64 {
	return int64(b.size)
}

func (b *base) HasNext() bool {
	return b.current < b.size
}

type sliceIterator struct {
	*base
	elements []types.T
}

func (s *sliceIterator) Next() types.T {
	e := s.elements[s.current]
	s.current++
	return e
}

type stringIt struct {
	*base
	elements []string
}

func (i *stringIt) Next() types.T {
	e := i.elements[i.current]
	i.current++
	return e
}

type sliceIt struct {
	*base
	sliceValue reflect.Value
}

func (s *sliceIt) Next() types.T {
	e := s.sliceValue.Index(s.current).Interface()
	s.current++
	return e
}

type supplierIt struct {
	get types.Supplier
}

func (s *supplierIt) GetSizeIfKnown() int64 {
	return unknownSize
}

func (s *supplierIt) HasNext() bool {
	return true
}

func (s *supplierIt) Next() types.T {
	return s.get()
}

type Sortable struct {
	List []types.T
	Cmp  types.Comparator
}

func (a *Sortable) Len() int {
	return len(a.List)
}

func (a *Sortable) Less(i, j int) bool {
	return a.Cmp(a.List[i], a.List[j]) < 0
}

func (a *Sortable) Swap(i, j int) {
	a.List[i], a.List[j] = a.List[j], a.List[i]
}

type endpoint interface {
	CompareTo(other endpoint) int
	Add(step int) endpoint
}

type epInt int

func (m epInt) CompareTo(other endpoint) int {
	return int(m - other.(epInt))
}

func (m epInt) Add(step int) endpoint {
	return m + epInt(step)
}

type epInt64 int64

func (m epInt64) CompareTo(other endpoint) int {
	return int(m - other.(epInt64))
}

func (m epInt64) Add(step int) endpoint {
	return m + epInt64(step)
}
