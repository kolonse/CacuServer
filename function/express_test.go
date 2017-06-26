package function

import (
	"testing"
)

func TestExpressParse(t *testing.T) {
	e := new(express)
	e.Parse("1 + 2 +  3+  4")
	t.Log(e.Call())

	e.Parse(" 1 * 2 +  3/  4 - 5 + 10 ")
	t.Log(e.Call())
}
