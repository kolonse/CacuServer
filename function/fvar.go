package function

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	MarkInt = iota
	MarkString
	MarkI
)

var (
	TypeMap  = make(map[string]interface{})
	NotValid = reflect.Value{}
)

type fvar struct {
	src    string
	k      string
	isleaf bool
	coos   []fvar
	args   []reflect.Value
	m      int
	i      int
}

// name[coo][coo][coo]
func (f *fvar) parse(index int, str string) (fvar, int) {
	var r fvar
	for i := index; i < len(str); {
		v := str[i]
		switch v {
		case '[':
			nv, ni := f.parse(i+1, str)
			r.coos = append(r.coos, nv)
			i = ni
		case ']':
			r.m = f.mark(r.k)
			r.isleaf = true
			if r.m == MarkInt {
				r.i, _ = strconv.Atoi(r.k)
			}
			if r.m == MarkString && len(r.coos) != 0 {
				r.isleaf = false
			}
			return r, i + 1
		default:
			r.k = string(append([]byte(r.k), v))
			i = i + 1
		}
	}
	r.m = f.mark(r.k)
	r.isleaf = true
	if r.m == MarkInt {
		r.i, _ = strconv.Atoi(r.k)
	}
	if r.m == MarkString && len(r.coos) != 0 {
		r.isleaf = false
	}
	return r, len(str)
}

func vprint(f *fvar, seq string) {
	if !f.isleaf {
		fmt.Printf("%v%v %v\n", seq, f.k, f.isleaf)
		for _, v := range f.coos {
			vprint(&v, seq+"\t")
		}
	} else {
		fmt.Printf("%v%v %v\n", seq, f.k, f.isleaf)
	}
}

// name[coo][coo][coo]
func (f *fvar) Parse(str string) {
	f.src = str
	v, _ := f.parse(0, str)
	f.k = v.k
	f.coos = v.coos
	f.m = v.m
	f.isleaf = v.isleaf
}

func (f *fvar) call(fv fvar, sarg reflect.Value, args ...reflect.Value) (reflect.Value, error) {
	var arg reflect.Value
	if fv.isleaf {
		arg = args[0]
	} else {
		arg = sarg
	}
	var rs reflect.Value
	switch fv.m {
	case MarkI: // 如果是 i 那么参数中必须传入 i 值
		i := int(args[1].Int())
		switch args[0].Type().Kind() {
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			rs = args[0].Index(i)
		case reflect.Interface:
			rs = reflect.ValueOf(args[0].Interface()).Index(i)
		default:
			panic(fmt.Errorf("i 下标规则 需要参数确保是 []interface{}"))
		}
	case MarkInt:
		switch args[0].Type().Kind() {
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			rs = args[0].Index(fv.i)
		case reflect.Map:
			rs = args[0].MapIndex(reflect.ValueOf(fv.k))
			if !rs.IsValid() {
				panic(fmt.Errorf("%v 变量值不存在", fv.k))
			}
		case reflect.Interface:
			ifc := args[0].Interface()
			ifcv := reflect.ValueOf(ifc)
			switch ifcv.Kind() {
			case reflect.Slice:
				fallthrough
			case reflect.Array:
				rs = ifcv.Index(fv.i)
			case reflect.Map:
				rs = ifcv.MapIndex(reflect.ValueOf(fv.k))
				if !rs.IsValid() {
					return NotValid, fmt.Errorf("%v 变量值不存在", fv.k)
				}
			default:
				return NotValid, fmt.Errorf("整数 下标规则 需要参数确保是 []interface{}/map[string]interface{}")
			}
		default:
			panic(fmt.Errorf("整数 下标规则 需要参数确保是 []interface{}/map[string]interface{}"))
		}
	case MarkString:
		switch arg.Type().Kind() {
		case reflect.Map:
			rs = arg.MapIndex(reflect.ValueOf(fv.k))
			if !rs.IsValid() {
				return NotValid, fmt.Errorf("%v 变量值不存在", fv.k)
			}
		case reflect.Interface:
			rs = reflect.ValueOf(arg.Interface()).MapIndex(reflect.ValueOf(fv.k))
			if !rs.IsValid() {
				return NotValid, fmt.Errorf("%v 变量值不存在", fv.k)
			}
		default:
			panic(fmt.Errorf("字符串 下标规则 需要参数确保是 map[string]interface{}"))
		}
	default:
		panic(fmt.Errorf("not support"))
	}
	if len(fv.coos) != 0 {
		for _, v := range fv.coos {
			if !v.isleaf {
				ns, err := f.call(v, sarg, args...)
				if err != nil {
					return NotValid, err
				}
				tp := ns.Type()
				kind := tp.Kind()
				if kind >= reflect.Int && kind <= reflect.Uint64 {
					v.i = int(ns.Int())
					v.m = MarkInt
					v.k = strconv.Itoa(v.i)
					v.coos = make([]fvar, 0)
					v.isleaf = true
				} else {
					v.m = MarkString
					v.k = fmt.Sprintf("%v", ns.Interface())
					v.coos = make([]fvar, 0)
					v.isleaf = true
				}
			}
			if len(args) > 1 {
				var nargs []reflect.Value
				nargs = append(nargs, rs)
				nargs = append(nargs, args[1:]...)
				var err error
				rs, err = f.call(v, sarg, nargs...)
				if err != nil {
					return NotValid, err
				}
			} else {
				var err error
				rs, err = f.call(v, sarg, rs)
				if err != nil {
					return NotValid, err
				}
			}
		}
	}
	if rs.Kind() == reflect.Interface {
		rs = reflect.ValueOf(rs.Interface())
	}
	return rs, nil
}

func (f *fvar) mark(v string) int {
	if v == "i" {
		return MarkI
	}
	for _, v1 := range v {
		if v1 > '9' || v1 < '0' {
			return MarkString
		}
	}

	return MarkInt
}

func (f *fvar) Call() (interface{}, error) {
	return f.call(*f, f.args[0], f.args...)
}

func (f *fvar) SetCallArgs(args ...interface{}) {
	f.args = make([]reflect.Value, 0)
	for _, v := range args {
		f.args = append(f.args, reflect.ValueOf(v))
	}
}

func NewFvar() *fvar {
	return new(fvar)
}
