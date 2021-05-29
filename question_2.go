package stream_test

import (
	"stream_test/stream"
	"stream_test/stream/types"
)

// - Q1: 计算一个 string 中小写字母的个数
func Question2Sub1(str string) int64 {
	return stream.OfSlice([]byte(str)).Filter(func(e types.T) bool {
		return e.(byte) >= 'a' && e.(byte) <= 'z'
	}).Count()
}

type Elem struct {
	Count int64
	S     string
}

// - Q2: 找出 []string 中，包含小写字母最多的字符串
func Question2Sub2(list []string) string {
	return stream.OfStrings(list...).Map(func(e types.T) types.R {
		str := e.(string)
		return Elem{
			Count: stream.OfSlice([]byte(str)).Filter(func(e types.T) bool {
				return e.(byte) >= 'a' && e.(byte) <= 'z'
			}).Count(),
			S: str,
		}
	}).Sorted(func(left types.T, right types.T) int {
		return int(left.(Elem).Count - right.(Elem).Count)
	}).Limit(1).FindFirst().Get().(Elem).S
}
