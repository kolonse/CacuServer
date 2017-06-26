/**
*	函数格式 ${TIME,exp,fmt}
*	exp为表达式,如果表达式中出现 day,hour,minute,second 等字段,会被替换成当前时间的对应的数字
*	例如 day - 1,假设当前时间为2017.6.14,那么会替换成 14 - 1,那么最后的时间就是2014.6.13
*	fmt 是时间格式化方式 YYYY,MM,DD,HH,mm,ss
 */
package function

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kolonse/CacuServer/lib"
)

var ftimeexp = regexp.MustCompile(`\$\{TIME,([^,\}]+)?,([^\}]+)}`)
var expcacumap = map[string]map[byte]func(int64) func() time.Time{
	"day": map[byte]func(int64) func() time.Time{
		'+': func(day int64) func() time.Time {
			return func() time.Time {
				return baseTime(day * 24 * int64(time.Hour))
			}
		},
		'-': func(day int64) func() time.Time {
			return func() time.Time {
				return baseTime(-day * 24 * int64(time.Hour))
			}
		},
	},
	"hour": map[byte]func(int64) func() time.Time{
		'+': func(hour int64) func() time.Time {
			return func() time.Time {
				return baseTime(hour * int64(time.Hour))
			}
		},
		'-': func(hour int64) func() time.Time {
			return func() time.Time {
				return baseTime(-hour * int64(time.Hour))
			}
		},
	},
	"minute": map[byte]func(int64) func() time.Time{
		'+': func(minute int64) func() time.Time {
			return func() time.Time {
				return baseTime(minute * int64(time.Minute))
			}
		},
		'-': func(minute int64) func() time.Time {
			return func() time.Time {
				return baseTime(-minute * int64(time.Minute))
			}
		},
	},
	"second": map[byte]func(int64) func() time.Time{
		'+': func(second int64) func() time.Time {
			return func() time.Time {
				return baseTime(second * int64(time.Second))
			}
		},
		'-': func(second int64) func() time.Time {
			return func() time.Time {
				return baseTime(-second * int64(time.Second))
			}
		},
	},
}

func baseTime(d int64) time.Time {
	return time.Now().Add(time.Duration(d))
}

type ftime struct {
	src    string
	result string
	format string
	fexp   func() time.Time
}

func (f *ftime) parseExp(exp string) {
	if len(exp) == 0 {
		f.fexp = func() time.Time {
			return baseTime(0)
		}
		return
	}
	var left string
	var mark byte
	s := 0
	for i := 0; i < len(exp); i++ {
		switch exp[i] {
		case '+':
			fallthrough
		case '-':
			left = strings.Trim(exp[0:i], " ")
			mark = exp[i]
			s = i + 1
		}
	}
	right, err := strconv.Atoi(strings.Trim(exp[s:], " "))
	if err != nil {
		panic(err)
	}
	f.fexp = expcacumap[left][mark](int64(right))
}

func (f *ftime) Parse(str string) {
	f.src = str
	matchs := ftimeexp.FindStringSubmatch(str)
	exp := strings.Trim(matchs[1], " ")
	f.parseExp(exp)
	f.format = strings.Trim(matchs[2], " ")
}

func (f *ftime) Call() (interface{}, error) {
	return lib.TimeFormat(f.fexp(), f.format), nil
}

func (f *ftime) SetCallArgs(args ...interface{}) {
}

func NewFtime() *ftime {
	return new(ftime)
}
