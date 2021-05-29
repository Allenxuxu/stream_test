package stream_test

import (
	"stream_test/stream"
	"stream_test/stream/types"
)

// - Q1: 输入一个整数 int，字符串string。将这个字符串重复n遍返回
func Question3Sub1(str string, n int) (ret string) {
	stream.RepeatN(str, int64(n)).ForEach(func(e types.T) {
		ret += e.(string)
	})

	return
}
