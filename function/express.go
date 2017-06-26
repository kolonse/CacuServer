package function

import (
	"strconv"
	"strings"

	"github.com/kolonse/CacuServer/lib"
)

var endFlag = byte('$')
var flagmap = map[byte]map[byte]int{
	'+': map[byte]int{'+': 0, '-': 0, '*': -1, '/': -1, '$': 1},
	'-': map[byte]int{'+': 0, '-': 0, '*': -1, '/': -1, '$': 1},
	'*': map[byte]int{'+': 1, '-': 1, '*': 0, '/': 0, '$': 1},
	'/': map[byte]int{'+': 1, '-': 1, '*': 0, '/': 0, '$': 1},
	'$': map[byte]int{'+': -1, '-': -1, '*': -1, '/': -1, '$': 0},
}

var expmap = map[byte]func(float64, float64) float64{
	'+': func(a, b float64) float64 { return a + b },
	'-': func(a, b float64) float64 { return a - b },
	'*': func(a, b float64) float64 { return a * b },
	'/': func(a, b float64) float64 { return a / b },
}

type express struct {
	result float64
}

func (e *express) cacu(snum *lib.Stack, smark *lib.Stack, flag byte) {
	if smark.Size() < 1 {
		smark.Push(flag)
		return
	}
	f := smark.Pop().(byte)
	if flagmap[f][flag] >= 0 {
		sv1 := snum.Pop().(float64)
		sv2 := snum.Pop().(float64)
		snum.Push(expmap[f](sv2, sv1))
		smark.Push(flag)
	} else {
		smark.Push(f)
		smark.Push(flag)
	}
	return
}

func (e *express) push(snum *lib.Stack, str string) {
	s := strings.Trim(str, " ")
	if len(s) == 0 {
		return
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	snum.Push(v)
}

func (e *express) Parse(str string) {
	s := 0
	snum := lib.NewStack()
	smark := lib.NewStack()
	for i := 0; i < len(str); i++ {
		switch str[i] {
		case '+':
			fallthrough
		case '-':
			fallthrough
		case '*':
			fallthrough
		case '/':
			e.push(snum, str[s:i])
			e.cacu(snum, smark, str[i])
			s = i + 1
		case ' ':
			e.push(snum, str[s:i])
			s = i + 1
		}
	}
	if s != len(str) {
		e.push(snum, str[s:])
	}
	e.cacu(snum, smark, endFlag)
	// 遍历栈空间对栈余留元素进行处理
	for !smark.Empty() {
		f := smark.Pop().(byte)
		if f == endFlag {
			continue
		}
		sv1 := snum.Pop().(float64)
		sv2 := snum.Pop().(float64)
		snum.Push(expmap[f](sv2, sv1))
	}
	e.result = snum.Pop().(float64)
}

func (e *express) Call() (interface{}, error) {
	return e.result, nil
}

func (e *express) SetCallArgs(args ...interface{}) {
}
