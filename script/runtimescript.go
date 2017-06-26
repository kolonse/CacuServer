/**
*	runtime 运行时脚本,需要在调用前传入全局参数 map[string][][]interface{},
*	和结果输出 map[string]interface{},以及当前计算数据位标 i
*	该脚本会运行逻辑,并有一定的输出结果会写到 map[string]interface{} 中
*	支持 求和 运算, 正则:${SUM,结果变量名,计算字段 Array[0][i][price],分片字段 Array[0][i][leader]}
*
 */
package script

import (
	"fmt"
	"regexp"

	"github.com/kolonse/CacuServer/function"
)

var runtimereg = regexp.MustCompile(`\$\{(SUM|FOR|LUASUM),[^\}]+}`)

//var runtimeregmap = map[string]*regexp.Regexp{
//	"SUM": regexp.MustCompile(`\$\{SUM,([^,\}]+),([^,\}]+),([^\}]+)?}`),
//}

var runtimefuncmap = map[string]func() function.Function{
	"SUM":    func() function.Function { return function.NewSum() },
	"FOR":    func() function.Function { return function.NewFfor() },
	"LUASUM": func() function.Function { return function.NewFluasum() },
}

type RuntimeScript struct {
	src  string
	args []interface{}
	f    function.Function
}

func (s *RuntimeScript) SetCallArgs(args ...interface{}) {
	s.args = make([]interface{}, 0)
	s.args = append(s.args, args...)
	s.f.SetCallArgs(s.args...)
}

func (s *RuntimeScript) Parse(str string) {
	s.src = str
	matchs := runtimereg.FindStringSubmatch(s.src)
	if len(matchs) < 2 {
		panic(fmt.Errorf("%v not support", str))
	}
	f, ok := runtimefuncmap[matchs[1]]
	if !ok {
		panic(fmt.Errorf("%v not support", str))
	}
	s.f = f()
	s.f.Parse(str)
}

func (s *RuntimeScript) Call() (interface{}, error) {
	return s.f.Call()
}

func NewRuntimeScript() *RuntimeScript {
	return new(RuntimeScript)
}
