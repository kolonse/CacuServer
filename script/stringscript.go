/**
*	string 脚本表示,最后直接输出字符串
 */
package script

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kolonse/CacuServer/function"
)

var stringregmap = map[string]*regexp.Regexp{
	"TIME": regexp.MustCompile(`\$\{TIME,([^,\}]+)?,([^\}]+)}`),
}

var stringfuncmap = map[string]func() function.Function{
	"TIME": func() function.Function { return function.NewFtime() },
}

type StringScript struct {
	src    string
	format string
	arr    []function.Function
}

func (s *StringScript) SetCallArgs(...interface{}) {

}

func (s *StringScript) Parse(str string) {
	s.src = str
	s.format = str
	bEnd := false
	for !bEnd {
		bEnd = true
		for key, reg := range stringregmap {
			match := reg.FindString(s.format)
			if len(match) == 0 {
				continue
			}
			s.format = strings.Replace(s.format, match, "%v", 1)
			f := stringfuncmap[key]()
			f.Parse(match)
			s.arr = append(s.arr, f)
			bEnd = false
		}
	}
}

func (s *StringScript) Call() (interface{}, error) {
	var args []interface{}
	for _, v := range s.arr {
		r, _ := v.Call()
		args = append(args, r)
	}
	return fmt.Sprintf(s.format, args...), nil
}

func NewStringScript() *StringScript {
	return new(StringScript)
}
