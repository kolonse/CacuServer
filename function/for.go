package function

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// ${FOR,OBJECT,proxyMoney,mysql $MysqlFcadmin  insert into proxyrecharge(proxyId,money) values(${KEY},${VALUE})}
var forreg = regexp.MustCompile(`^\$\{FOR,(OBJECT),([^,\}]+),(.+)}$`)

const (
	ForObject    = "OBJECT"
	ForKeyFlag   = "${KEY}"
	ForValueFlag = "${VALUE}"
)

type ffor struct {
	src   string
	args  []reflect.Value
	vtype string
	vname string
	dostr string
}

func (f *ffor) Parse(str string) {
	f.src = str
	matchs := forreg.FindStringSubmatch(str)
	f.vtype = matchs[1]
	f.vname = matchs[2]
	f.dostr = matchs[3]
	if f.vtype == ForObject {
		f.dostr = strings.Replace(f.dostr, ForKeyFlag, "%[1]v", -1)
		f.dostr = strings.Replace(f.dostr, ForValueFlag, "%[2]v", -1)
	}
}

func (f *ffor) callObject() (interface{}, error) {
	arg := f.args[0]
	if arg.Kind() != reflect.Map {
		return nil, fmt.Errorf("确保参数传递为 map")
	}
	obj := arg.MapIndex(reflect.ValueOf(f.vname))
	if obj.Kind() == reflect.Interface {
		obj = reflect.ValueOf(obj.Interface())
	}
	if !obj.IsValid() {
		return nil, fmt.Errorf("%v 对应的存储值不存在", f.vname)
	}
	if obj.Kind() != reflect.Map {
		return nil, fmt.Errorf("OBJECT 对象,确保参数传递为 map")
	}
	keys := obj.MapKeys()
	return func(cb func(str string)) {
		for _, k := range keys {
			v := obj.MapIndex(k)
			cb(fmt.Sprintf(f.dostr, k.Interface(), v.Interface()))
		}
	}, nil
}

func (f *ffor) Call() (interface{}, error) {
	if f.vtype == ForObject {
		return f.callObject()
	}
	return nil, fmt.Errorf("not support")
}

func (f *ffor) SetCallArgs(args ...interface{}) {
	f.args = make([]reflect.Value, 0)
	for _, v := range args {
		f.args = append(f.args, reflect.ValueOf(v))
	}
}

func NewFfor() *ffor {
	return new(ffor)
}
