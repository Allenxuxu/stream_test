package stream

import (
	"errors"
	"reflect"

	"stream_test/stream/optional"
	"stream_test/stream/types"
)

var (
	// ErrNotSlice a error to panic when call Slice but argument is not slice
	ErrNotSlice = errors.New("not slice")
)

func Of(elements ...types.T) Stream {
	return newHead(it(elements...))
}

func OfStrings(element ...string) Stream {
	return newHead(&stringIt{
		base: &base{
			current: 0,
			size:    len(element),
		},
		elements: element,
	})
}

// OfSlice return a Stream. the input parameter `slice` must be a slice.
// if input is nil, return a empty Stream( same as Of() )
func OfSlice(slice types.T) Stream {
	if optional.IsNil(slice) {
		return Of()
	}
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		panic(ErrNotSlice)
	}
	value := reflect.ValueOf(slice)
	it := &sliceIt{
		base: &base{
			current: 0,
			size:    value.Len(),
		},
		sliceValue: value,
	}
	return newHead(it)
}

func Repeat(e types.T) Stream {
	return newHead(withSupplier(func() types.T {
		return e
	}))
}

func RepeatN(e types.T, count int64) Stream {
	return Repeat(e).Limit(count)
}
