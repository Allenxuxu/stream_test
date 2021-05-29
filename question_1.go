package stream_test

import (
	"stream_test/stream"
	"stream_test/stream/types"
)

// - Q1: 输入 employees，返回 年龄 >22岁 的所有员工，年龄总和
func Question1Sub1(employees []*Employee) int64 {
	return stream.OfSlice(employees).Filter(func(e types.T) bool {
		return *e.(*Employee).Age > 22
	}).Count()
}

// - Q2: - 输入 employees，返回 id 最小的十个员工，按 id 升序排序
func Question1Sub2(employees []*Employee) (ret []*Employee) {
	stream.OfSlice(employees).Sorted(func(left types.T, right types.T) int {
		return int(left.(*Employee).Id - right.(*Employee).Id)
	}).Limit(10).ForEach(func(e types.T) {
		ret = append(ret, e.(*Employee))
	})

	return ret
}

// - Q3: - 输入 employees，对于没有手机号为0的数据，随机填写一个
func Question1Sub3(employees []*Employee) []*Employee {
	stream.OfSlice(employees).Filter(func(e types.T) bool {
		return *e.(*Employee).Phone == ""
	}).ForEach(func(e types.T) {
		e.(*Employee).Phone = randomPhone()
	})
	return employees
}

// - Q4: - 输入 employees ，返回一个map[int][]int，其中 key 为 员工年龄 Age，value 为该年龄段员工ID
func Question1Sub4(employees []*Employee) map[int][]int64 {
	ret := make(map[int][]int64)
	stream.OfSlice(employees).ForEach(func(e types.T) {
		employee := e.(*Employee)
		ret[*employee.Age] = append(ret[*employee.Age], employee.Id)
	})
	return ret
}
