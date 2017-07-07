package function

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/kolonse/CacuServer/lib"
	lua "github.com/yuin/gopher-lua"
)

// ${LUASUM,lua file,var name,div name,arg1,arg2...}
var (
	luareg = regexp.MustCompile(`^\$\{LUASUM,([^,]+),([^,]+),([^,]+),(.+)}$`)
	luamap = lib.NewSafeMap()
)

type fluasum struct {
	src      string
	args     []interface{}
	luafile  string // lua 文件
	vname    string // 变量名
	dname    *fvar  // 拆分名字
	cargs    []*fvar
	luaState *lua.LState
}

func (f *fluasum) Parse(str string) {
	f.src = str
	matchs := luareg.FindStringSubmatch(str)
	f.luafile = strings.Trim(matchs[1], " ")
	f.vname = strings.Trim(matchs[2], " ")
	fv := NewFvar()
	fv.Parse(strings.Trim(matchs[3], " "))
	f.dname = fv
	args := strings.Split(strings.Trim(matchs[4], " "), ",")
	for _, v := range args {
		v := strings.Trim(v, " ")
		if len(v) == 0 {
			continue
		}
		fv = NewFvar()
		fv.Parse(v)
		f.cargs = append(f.cargs, fv)
	}
	l, ok := luamap.MapIndex(f.luafile)
	if !ok {
		l = lua.NewState()
		if err := l.(*lua.LState).DoFile(f.luafile); err != nil {
			panic(err)
		}
		f.luaState = l.(*lua.LState)
		luamap.SetMapIndex(f.luafile, l)
	} else {
		f.luaState = l.(*lua.LState)
	}
}

func (f *fluasum) Call() (interface{}, error) {
	iparam := reflect.ValueOf(f.args[0])
	tname := reflect.ValueOf(f.vname)
	// 读取值 vname 字段
	r := iparam.MapIndex(tname)
	if !r.IsValid() {
		r = reflect.MakeMap(divmap)
	}
	if r.Kind() == reflect.Interface {
		r = reflect.ValueOf(r.Interface())
	}

	// 获取拆分字段
	f.dname.SetCallArgs(f.args...)
	di, err := f.dname.Call()
	if err != nil {
		return nil, err
	}
	d := di.(reflect.Value)
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
	// 调用lua接口计算结果v1
	var luaArgs []lua.LValue
	for _, v := range f.cargs {
		v.SetCallArgs(f.args...)
		di, e := v.Call()
		if e != nil {
			//			return nil, e
			di = reflect.ValueOf("_ilegal_")
		}
		d := di.(reflect.Value)
		dk := d.Kind()
		if dk == reflect.Bool {
			luaArgs = append(luaArgs, lua.LBool(d.Bool()))
		} else if dk >= reflect.Int && dk <= reflect.Uint64 {
			luaArgs = append(luaArgs, lua.LNumber(d.Int()))
		} else if dk >= reflect.Float32 && dk <= reflect.Float64 {
			luaArgs = append(luaArgs, lua.LNumber(d.Float()))
		} else if dk == reflect.String {
			luaArgs = append(luaArgs, lua.LString(d.String()))
		} else {
			return nil, fmt.Errorf("%v not support", dk.String())
		}
	}
	if err := f.luaState.CallByParam(lua.P{
		Fn:      f.luaState.GetGlobal("entry"),
		NRet:    1,
		Protect: true,
	}, luaArgs...); err != nil {
		return nil, err
	}
	vl := f.luaState.Get(-1) // returned value
	f.luaState.Pop(1)        // remove received value
	if vl.Type() != lua.LTNumber {
		return nil, fmt.Errorf("%v lua 脚本返回值非预期数值类型")
	}
	v2, _ := strconv.ParseFloat(vl.String(), 64)
	/////
	sumv = reflect.ValueOf(sumv.Float() + v2)
	r.SetMapIndex(rtv, sumv)
	iparam.SetMapIndex(tname, r)
	return iparam.Interface(), nil
}

func (f *fluasum) SetCallArgs(args ...interface{}) {
	f.args = args
}

func NewFluasum() *fluasum {
	return new(fluasum)
}
