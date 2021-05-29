package stream

import (
	"stream_test/stream/optional"
	"stream_test/stream/types"
)

type Stream interface {
	Filter(types.Predicate) Stream
	Map(types.Function) Stream

	Sorted(types.Comparator) Stream
	Limit(int64) Stream

	ForEach(types.Consumer)
	FindFirst() optional.Optional
	Count() int64

	// ...
}
