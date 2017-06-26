package function

import (
	"testing"
)

func TestSum(t *testing.T) {
	e := new(sum)
	m := make(map[string]interface{})
	Array := make([]interface{}, 1)
	Array[0] = make([]interface{}, 5)
	for i, _ := range Array[0].([]interface{}) {
		Array[0].([]interface{})[i] = make(map[string]interface{})
		Array[0].([]interface{})[i].(map[string]interface{})["price"] = 10
		Array[0].([]interface{})[i].(map[string]interface{})["leader"] = 0
	}
	m["Array"] = Array
	e.Parse("${SUM,leaderMoney,Array[0][i][price],Array[0][i][leader]}")
	for i, _ := range Array[0].([]interface{}) {
		e.SetCallArgs(m, i)
		mi, err := e.Call()
		if err != nil {
			t.Error(err)
		}
		m = mi.(map[string]interface{})
	}
	if m["leaderMoney"].(map[string]float64)["0"] != 50 {
		t.Error(m, "结果不匹配")
	}
}
