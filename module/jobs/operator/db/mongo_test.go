package db

import (
	"testing"
)

func TestMgo(t *testing.T) {
	m, e := NewMMgo(`mongodb://120.27.196.178/storeServer carduseds _sel_{"_id":0,"data.leader":1,"data.unionId":1,"data.cardcount":1}_ect_ {"data.dateAt":{"$gte":"${TIME,day - 1,YYYY-01-DD 00:00:00}","$lte":"${TIME,day - 1,YYYY-MM-DD 23:59:59}"}}`)
	if e != nil {
		t.Error(e)
		return
	}
	ret, err := m.Count()
	t.Log(err)
	t.Log(ret)

	r, _ := m.Read()
	t.Log(r)
}
