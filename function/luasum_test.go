package function

import (
	"testing"
)

func TestLuaSum1(t *testing.T) {
	e := new(flua)
	m := make(map[string]interface{})
	m["utype"] = 0
	m["usharingrate"] = 10
	m["psharingrate"] = 10
	m["ursharingrate"] = 10
	m["price"] = 100000
	m["leader"] = "10080"
	e.Parse(`${LUASUM,../conf/leaderIncome.lua,leaderIncome,leader,utype,usharingrate,psharingrate,ursharingrate,price}`)
	e.SetCallArgs(m)
	m1, err := e.Call()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(m1.(map[string]interface{})["leaderIncome"].(map[string]float64)["10080"])
	m["utype"] = 1
	m1, err = e.Call()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(m1.(map[string]interface{})["leaderIncome"].(map[string]float64)["10080"])
}

func TestLuaSum2(t *testing.T) {
	e := new(flua)
	m := make(map[string]interface{})
	m["utype"] = 1
	m["usharingrate"] = 10
	m["psharingrate"] = 10
	m["ursharingrate"] = 10
	m["price"] = 100000
	m["leader"] = "10080"
	e.Parse(`${LUASUM,../conf/leaderIncome.lua,leaderIncome,leader,utype,usharingrate,psharingrate,ursharingrate,price}`)
	e.SetCallArgs(m)
	m1, err := e.Call()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(m1.(map[string]interface{})["leaderIncome"].(map[string]float64)["10080"])
}
