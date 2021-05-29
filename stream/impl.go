package stream

import (
	"sort"

	"stream_test/stream/optional"
	"stream_test/stream/types"
)

type stream struct {
	source iterator
	prev   *stream
	wrap   func(stage) stage
}

func newHead(source iterator) *stream {
	return &stream{source: source}
}

func newNode(prev *stream, wrap func(down stage) stage) *stream {
	return &stream{
		source: prev.source,
		prev:   prev,
		wrap:   wrap,
	}
}

func (s *stream) terminal(ts *terminalStage) {
	stage := s.wrapStage(ts)
	source := s.source
	stage.Begin(source.GetSizeIfKnown())
	for source.HasNext() && !stage.CanFinish() {
		stage.Accept(source.Next())
	}
	stage.End()
}

func (s *stream) wrapStage(terminalStage stage) stage {
	stage := terminalStage
	for i := s; i.prev != nil; i = i.prev {
		stage = i.wrap(stage)
	}
	return stage
}

func (s *stream) Filter(test types.Predicate) Stream {
	return newNode(s, func(down stage) stage {
		return newChainedStage(down, begin(func(int64) {
			down.Begin(unknownSize) // 过滤后个数不确定
		}), action(func(t types.T) {
			if test(t) {
				down.Accept(t)
			}
		}))
	})
}

func (s *stream) Map(apply types.Function) Stream {
	return newNode(s, func(down stage) stage {
		return newChainedStage(down, action(func(t types.T) {
			down.Accept(apply(t))
		}))
	})
}

func (s *stream) Sorted(cmp types.Comparator) Stream {
	return newNode(s, func(down stage) stage {
		var list []types.T
		return newChainedStage(down, begin(func(size int64) {
			if size > 0 {
				list = make([]types.T, 0, size)
			} else {
				list = make([]types.T, 0)
			}
			down.Begin(size)
		}), action(func(t types.T) {
			list = append(list, t)
		}), end(func() {
			a := &Sortable{
				List: list,
				Cmp:  cmp,
			}
			sort.Sort(a)
			down.Begin(int64(len(a.List)))
			i := it(a.List...)
			for i.HasNext() && !down.CanFinish() {
				down.Accept(i.Next())
			}
			list = nil
			a = nil
			down.End()
		}))
	})
}

func (s *stream) Limit(maxSize int64) Stream {
	return newNode(s, func(down stage) stage {
		count := int64(0)
		return newChainedStage(down, begin(func(size int64) {
			if size > 0 {
				if size > maxSize {
					size = maxSize
				}
			}
			down.Begin(size)
		}), action(func(t types.T) {
			if count < maxSize {
				down.Accept(t)
			}
			count++
		}), canFinish(func() bool {
			return count == maxSize
		}))
	})
}

func (s *stream) ForEach(consumer types.Consumer) {
	s.terminal(newTerminalStage(consumer))
}

func (s *stream) ReduceWith(initValue types.R, accumulator func(types.R, types.T) types.R) types.R {
	var result = initValue
	s.terminal(newTerminalStage(func(t types.T) {
		result = accumulator(result, t)
	}))
	return result
}

func (s *stream) FindFirst() optional.Optional {
	var result types.T
	var find = false
	s.terminal(newTerminalStage(func(t types.T) {
		if !find {
			result = t
			find = true
		}
	}, canFinish(func() bool {
		return find
	})))
	return optional.OfNullable(result)
}

func (s *stream) Count() int64 {
	return s.ReduceWith(int64(0), func(count types.R, t types.T) types.R {
		return count.(int64) + 1
	}).(int64)
}
