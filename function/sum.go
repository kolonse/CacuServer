package function

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

type sum struct {
	src     string
	args    []interface{}
	name    string
	sumflag *fvar
	sumdiv  *fvar
}

var sumexp = regexp.MustCompile(`\$\{SUM,([^,\}]+),([^,\}]+),([^\}]+)}`)
var divmap = reflect.TypeOf(map[string]float64{})

func (f *sum) Parse(str string) {
	f.src = str
	matchs := sumexp.FindStringSubmatch(str)
	f.name = matchs[1]
	f.sumflag = NewFvar()
	f.sumflag.Parse(matchs[2])
	f.sumdiv = NewFvar()
	f.sumdiv.Parse(matchs[3])
}

func (f *sum) Call() (interface{}, error) {
	f.sumflag.SetCallArgs(f.args...)
	vi, err := f.sumflag.Call()
	if err != nil {
		return nil, err
	}
	v := vi.(reflect.Value)
	iparam := reflect.ValueOf(f.args[0])
	vkd := v.Kind()
	vl := float64(0)
	if vkd >= reflect.Int && vkd <= reflect.Uint64 {
		vl = float64(v.Int())
	} else if vkd == reflect.Float32 || vkd == reflect.Float64 {
		vl = v.Float()
	} else if vkd == reflect.String {
		v2, err := strconv.ParseFloat(v.String(), 64)
		if err != nil {
			v2 = 0
		}
		vl = v2
	}
	// 读取拆分字段
	f.sumdiv.SetCallArgs(f.args...)
	di, err := f.sumdiv.Call()
	if err != nil {
		return nil, err
	}
	d := di.(reflect.Value)
	// 根据变量名字 获取统计 拆分字段 数值map
	tname := reflect.ValueOf(f.name)
	r := iparam.MapIndex(tname)
	if !r.IsValid() {
		r = reflect.MakeMap(divmap)
	}
	if r.Kind() == reflect.Interface {
		r = reflect.ValueOf(r.Interface())
	}
	rt := ""
	kd := d.Kind()
	if kd == reflect.Bool {
		rt = fmt.Sprintf("%v", d.Bool())
	} else if kd >= reflect.Int && kd <= reflect.Uint64 {
		rt = fmt.Sprintf("%v", d.Int())
	} else if kd == reflect.Float32 || kd == reflect.Float64 {
		rt = fmt.Sprintf("%v", d.Float())
	} else if kd == reflect.String {
		rt = fmt.Sprintf("%v", d.String())
	} else {
		return nil, fmt.Errorf("分区字段必须是字符串或者数字")
	}
	rtv := reflect.ValueOf(rt)
	sumv := r.MapIndex(rtv)
	if !sumv.IsValid() {
		sumv = reflect.ValueOf(float64(0))
	}
	sumv = reflect.ValueOf(sumv.Float() + vl)
	r.SetMapIndex(rtv, sumv)
	iparam.SetMapIndex(tname, r)
	return iparam.Interface(), nil
}

func (f *sum) SetCallArgs(args ...interface{}) {
	f.args = make([]interface{}, 0)
	f.args = append(f.args, args...)
}

func NewSum() *sum {
	return new(sum)
}
