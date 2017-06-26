package script

import (
	"testing"
)

func TestRuntimeScript(t *testing.T) {
	script := NewRuntimeScript()
	script.Parse(`${SUM,leaderMoney,Array[0][i][price],Array[0][i][leader]}`)

	m := make(map[string]interface{})
	Array := make([]interface{}, 1)
	Array[0] = make([]interface{}, 5)
	for i, _ := range Array[0].([]interface{}) {
		Array[0].([]interface{})[i] = make(map[string]interface{})
		Array[0].([]interface{})[i].(map[string]interface{})["price"] = 10
		Array[0].([]interface{})[i].(map[string]interface{})["leader"] = 0
	}
	m["Array"] = Array
	for i, _ := range Array[0].([]interface{}) {
		script.SetCallArgs(m, i)
		mi, err := script.Call()
		if err != nil {
			t.Error(err)
		}
		m = mi.(map[string]interface{})
	}
	if m["leaderMoney"].(map[string]float64)["0"] != 50 {
		t.Error(m, "结果不匹配")
	}
}
