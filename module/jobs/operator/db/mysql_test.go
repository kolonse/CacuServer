package db

import (
	"testing"
)

func TestMysql(t *testing.T) {
	m, _ := NewMMysql("admin:123456@tcp(120.27.196.178:3306)/fcpay?charset=utf8 select count(*) from pay where time_end >='${TIME,day - 1,YYYY05DD000000}' &&  time_end <='${TIME,day - 1,YYYYMMDD235959}'")
	ret, err := m.Count()
	t.Log(err)
	t.Log(ret)
}

func TestMysqlWrite(t *testing.T) {
	m, _ := NewMMysql("admin:123456@tcp(120.27.196.178:3306)/fcadmin?charset=utf8 insert into leadercarduse(leader,date,acount,fcount,rcount,scount,dcount) values(10086,'${TIME,day - 1,YYYY-MM-DD}',0,0,0,0,0)")
	err := m.Write()
	t.Log(err)
}
