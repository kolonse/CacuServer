package script

import (
	"testing"
)

func TestStringScript(t *testing.T) {
	script := new(StringScript)
	script.Parse("select count(*) from pay where time_end >=${TIME,day - 1,YYYY-MM-DD 00:00:00} &&  time_end <=${TIME,day - 1,YYYY-MM-DD 23:59:59}")
	t.Log(script.Call())
}
