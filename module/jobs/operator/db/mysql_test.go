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

//func TestMysqlWrite(t *testing.T) {
//	m, _ := NewMMysql("admin:123456@tcp(120.27.196.178:3306)/fcadmin?charset=utf8 INSERT INTO proxyrecharge(createdAt,updatedAt) VALUES('2017-06-18 23:59:02','2017-06-18 23:59:02')")
//	err := m.Write()
//	t.Log(err)
//}
