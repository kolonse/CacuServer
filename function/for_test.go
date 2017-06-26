package function

import (
	"testing"
)

func TestFor(t *testing.T) {
	e := new(ffor)
	e.Parse(`${FOR,OBJECT,proxyMoney,mysql admin:123456@tcp(120.27.196.178:3306)/fcadmin?charset=utf8  insert into proxyrecharge(proxyId,money,createdAt,updatedAt) values(${KEY},${VALUE},'${TIME,day - 0,YYYY-MM-DD HH:mm:ss}','${TIME,day - 0,YYYY-MM-DD HH:mm:ss}')}`)
	m := make(map[string]interface{})
	m1 := make(map[string]float64)
	m1["34"] = float64(10.0)
	m["proxyMoney"] = m1
	e.SetCallArgs(m)
	f, err := e.Call()
	if err != nil {
		t.Error(err)
		return
	}
	f.(func(func(string)))(func(str string) {
		t.Log(str)
	})
}
