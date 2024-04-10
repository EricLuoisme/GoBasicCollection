package stack

type Stack struct {
	data []string // 使用slice
}

func (s *Stack) Push(x string) {
	s.data = append(s.data, x)
}

func (s *Stack) Pop() string {
	// 取出最后一位
	n := len(s.data) - 1
	res := s.data[n]
	// 将array置空
	s.data[n] = ""
	// 使用slice的方式进行更新
	s.data = s.data[:n]
	return res
}

func (s *Stack) Size() int {
	return len(s.data)
}
