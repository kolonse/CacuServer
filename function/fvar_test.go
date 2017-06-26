package function

import (
	"reflect"
	"testing"
)

func TestFVar(t *testing.T) {
	v := NewFvar()
	v.Parse("test[1]")
	m := make(map[string]interface{})
	m["test"] = make(map[string]int)
	m["test"].(map[string]int)["1"] = 10
	m["test"].(map[string]int)["2"] = 11
	v.SetCallArgs(m)
	if i, err := v.Call(); err != nil && i.(reflect.Value).Int() != 10 {
		t.Error("not pass")
	}
}

func TestFVar2(t *testing.T) {
	v := NewFvar()
	v.Parse("test[1]")
	m := make(map[string]interface{})
	m["test"] = make([]interface{}, 2)
	m["test"].([]interface{})[0] = 1
	m["test"].([]interface{})[1] = 2
	v.SetCallArgs(m)
	if i, err := v.Call(); err != nil && i.(reflect.Value).Int() != 2 {
		t.Error("not pass")
	}
}

func TestFVar3(t *testing.T) {
	v := NewFvar()
	v.Parse("test[i]")
	m := make(map[string]interface{})
	m["test"] = make([]interface{}, 2)
	m["test"].([]interface{})[0] = 1
	m["test"].([]interface{})[1] = 2
	v.SetCallArgs(m, 1)
	if i, err := v.Call(); err != nil && i.(reflect.Value).Int() != 2 {
		t.Error("not pass")
	}
}

func TestFVar4(t *testing.T) {
	v := NewFvar()
	v.Parse("test[test2[i]]")
	m := make(map[string]interface{})
	m["test"] = make([]interface{}, 2)
	m["test"].([]interface{})[0] = 3
	m["test"].([]interface{})[1] = 4
	m["test2"] = make([]interface{}, 2)
	m["test2"].([]interface{})[0] = 0
	m["test2"].([]interface{})[1] = 1
	v.SetCallArgs(m, 1)
	if i, err := v.Call(); err != nil && i.(reflect.Value).Int() != 4 {
		t.Error("not pass")
	}
}

func TestFVar5(t *testing.T) {
	v := NewFvar()
	v.Parse("test[test2[i]]")
	m := make(map[string]interface{})
	m["test"] = make(map[string]interface{})
	m["test"].(map[string]interface{})["1"] = "tt1"
	m["test"].(map[string]interface{})["2"] = "tt2"
	m["test2"] = make([]interface{}, 2)
	m["test2"].([]interface{})[0] = 0
	m["test2"].([]interface{})[1] = 1
	v.SetCallArgs(m, 1)
	if i, err := v.Call(); err != nil && i.(reflect.Value).String() != "tt1" {
		t.Error("not pass")
	}
}

func TestFVar6(t *testing.T) {
	v := NewFvar()
	v.Parse("GObject[1][Array[0][0][labourunionId]][proxyId]")
	m := make(map[string]interface{})
	o := make([]interface{}, 2)
	o1 := make(map[string]interface{})
	o2 := make(map[string]string)
	o2["proxyId"] = "10"
	o1["33"] = o2
	o[1] = o1
	m["GObject"] = o

	a := make([]interface{}, 1)
	a1 := make([]interface{}, 1)
	a2 := make(map[string]string)
	a2["labourunionId"] = "33"
	a1[0] = a2
	a[0] = a1
	m["Array"] = a
	v.SetCallArgs(m)
	if i, err := v.Call(); err != nil && i.(reflect.Value).String() != "10" {
		t.Error("not pass")
	}
}
