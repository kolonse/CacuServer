package lib

type Stack []interface{}

func (s *Stack) Push(e interface{}) {
	//	switch e.(type) {
	//	case float64:
	//		println("push", e.(float64))
	//	}

	*s = append(*s, e)
}

func (s *Stack) Pop() interface{} {
	if s.Empty() {
		return nil
	}
	r := (*s)[s.Size()-1]
	//	switch r.(type) {
	//	case float64:
	//		println("pop", r.(float64))
	//	}

	*s = (*s)[0 : s.Size()-1]
	return r
}

func (s *Stack) Empty() bool {
	if len(*s) == 0 {
		return true
	}
	return false
}

func (s *Stack) Size() int {
	return len(*s)
}
func NewStack() *Stack {
	return new(Stack)
}
