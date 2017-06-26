package function

import (
	"testing"
)

func TestFtimeParse(t *testing.T) {
	e := new(ftime)
	e.Parse("${TIME,day-1,YYYY-MM-DD HH:mm:ss}")
	t.Log(e.Call())

	e.Parse("${TIME,,YYYY-MM-DD HH:mm:ss}")
	t.Log(e.Call())
}
