package db

import (
	"fmt"
	"strconv"
	//	"reflect"
	"strings"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kolonse/CacuServer/conf"
	//	"github.com/kolonse/CacuServer/lib"
	"github.com/kolonse/CacuServer/script"
)

//object mark
const (
	by = "BY_"
)

var (
//	mysqlpool = lib.NewSafeMap()
)

type mmysql struct {
	connStr string
	sql     string
	sqtSql  script.Script
	offset  int
	kObj    string
	isend   bool
}

func (p *mmysql) init(param string) error {
	p.offset = 0
	index := strings.Index(param, " ")
	if index == -1 {
		return fmt.Errorf("[%v] db string format error", param)
	}
	p.connStr = strings.Trim(param[0:index], " ")
	p.sql = strings.Trim(param[index+1:], " ")
	if p.sql[0:len(by)] == by {
		index = strings.Index(p.sql, " ")
		if index == -1 {
			return fmt.Errorf("[%v] db string format error", param)
		}
		p.kObj = p.sql[len(by):index]
		p.sql = p.sql[index+1:]
	}
	index = strings.Index(p.sql, "limit")
	if index != -1 {
		return fmt.Errorf("[%v] sql not allow limit", param)
	}
	p.sqtSql = script.NewStringScript()
	p.sqtSql.Parse(p.sql)
	//	if _, ok := mysqlpool.MapIndex(p.connStr); !ok {
	//		db, err := sql.Open("mysql", p.connStr)
	//		if err != nil {
	//			return err
	//		}
	//		mysqlpool.SetMapIndex(p.connStr, db)
	//	}
	return nil
}

func (p *mmysql) Count() (int, error) {
	db, err := sql.Open("mysql", p.connStr)
	if err != nil {
		return 0, err
	}
	si, _ := p.sqtSql.Call()
	fmt.Println("执行sql:", si.(string))
	rows, err := db.Query(si.(string))
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		rows.Scan(scanArgs...)
		break
	}
	value := string(values[0])
	return strconv.Atoi(value)
}
func (p *mmysql) Reset() {
	p.offset = 0
	p.isend = false
}

func (p *mmysql) Read() (interface{}, error) {
	if p.isend {
		return nil, nil
	}
	db, err := sql.Open("mysql", p.connStr)
	if err != nil {
		return nil, err
	}
	si, _ := p.sqtSql.Call()
	sqlStr := si.(string) +
		" limit " +
		strconv.Itoa(p.offset) +
		"," +
		strconv.Itoa(conf.ReadCountLimit)
	fmt.Println("执行sql:", sqlStr)
	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	var ret interface{}
	if len(p.kObj) == 0 {
		ret = make([]interface{}, 0)
	} else {
		ret = make(map[string]interface{})
	}
	for i := range values {
		scanArgs[i] = &values[i]
	}
	count := 0
	for rows.Next() {
		rows.Scan(scanArgs...)
		var value string
		obj := make(map[string]string)
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "null"
			} else {
				value = string(col)
			}
			obj[columns[i]] = value
		}
		if len(p.kObj) == 0 {
			ret = append(ret.([]interface{}), obj)
		} else {
			ret.(map[string]interface{})[obj[p.kObj]] = obj
		}
		count = count + 1
	}
	p.offset = p.offset + count
	if count < conf.ReadCountLimit {
		p.isend = true
	}
	return ret, nil
}

func (p *mmysql) Write() error {
	db, err := sql.Open("mysql", p.connStr)
	if err != nil {
		return err
	}
	si, _ := p.sqtSql.Call()
	fmt.Println("执行sql:", si.(string))
	_, err = db.Exec(si.(string))
	if err != nil {
		return err
	}
	return nil
}

func NewMMysql(param string) (*mmysql, error) {
	r := new(mmysql)
	e := r.init(param)
	return r, e
}
