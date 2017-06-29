package db

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kolonse/CacuServer/conf"
	//	"github.com/kolonse/CacuServer/lib"
	"github.com/kolonse/CacuServer/script"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//object mark
const (
	sel_begin = "_sel_"
	sel_end   = "_ect_"
)

var (
//	mongopool = lib.NewSafeMap()
)

type mmgo struct {
	connStr string
	collStr string
	sql     string
	slt     string
	sltBson bson.M
	sqtSql  script.Script
	offset  int
	kObj    string
	isend   bool
}

func (p *mmgo) init(param string) error {
	p.offset = 0
	index := strings.Index(param, " ")
	if index == -1 {
		return fmt.Errorf("[%v] db string format error", param)
	}
	p.connStr = strings.Trim(param[0:index], " ")
	index2 := strings.Index(param[index+1:], " ")
	if index2 == -1 {
		return fmt.Errorf("[%v] db string format error", param)
	}
	p.collStr = strings.Trim(param[index+1:index+1+index2], " ")
	p.sql = strings.Trim(param[index+index2+1:], " ")
	// mongo select 字段内容提取
	if p.sql[0:len(sel_begin)] == sel_begin {
		index = strings.Index(p.sql, sel_end)
		if index == -1 {
			return fmt.Errorf("[%v] db string format error", param)
		}
		p.slt = strings.Trim(p.sql[len(sel_begin):index], " ")
		p.sql = strings.Trim(p.sql[index+len(sel_end):], " ")
		var b bson.M
		err := json.Unmarshal([]byte(p.slt), &b)
		if err != nil {
			return err
		}
		p.sltBson = b
	}
	p.sqtSql = script.NewStringScript()
	p.sqtSql.Parse(p.sql)
	//	if _, ok := mysqlpool.MapIndex(p.connStr); !ok {
	//		sess, err := mgo.Dial(p.connStr)
	//		if err != nil {
	//			return err
	//		}
	//		mongopool.SetMapIndex(p.connStr, sess.DB("").C(p.collStr))
	//	}
	return nil
}

func (p *mmgo) Count() (int, error) {
	sess, err := mgo.Dial(p.connStr)
	if err != nil {
		return 0, err
	}
	c := sess.DB("").C(p.collStr)
	s, err := p.sqtSql.Call()
	if err != nil {
		return 0, err
	}
	var b bson.M
	err = json.Unmarshal([]byte(s.(string)), &b)
	if err != nil {
		return 0, err
	}
	fmt.Println("执行 mongo count:", s.(string))
	return c.Find(b).Count()
}

func (p *mmgo) Reset() {
	p.offset = 0
	p.isend = false
}

func (p *mmgo) Read() (interface{}, error) {
	if p.isend {
		return nil, nil
	}
	sess, err := mgo.Dial(p.connStr)
	if err != nil {
		return nil, err
	}
	c := sess.DB("").C(p.collStr)
	s, err := p.sqtSql.Call()
	if err != nil {
		return nil, err
	}
	var b bson.M
	err = json.Unmarshal([]byte(s.(string)), &b)
	if err != nil {
		return nil, err
	}
	fmt.Println("执行 mongo read:", s.(string))
	var datas []interface{}
	if len(p.slt) != 0 {
		err = c.Find(b).Select(p.sltBson).Skip(p.offset).Limit(conf.ReadCountLimit).All(&datas)
	} else {
		err = c.Find(b).Skip(p.offset).Limit(conf.ReadCountLimit).All(&datas)
	}
	if err != nil {
		return nil, err
	}
	if len(datas) < conf.ReadCountLimit {
		p.isend = true
	}
	return datas, nil
}

func (p *mmgo) Write() error {
	// mongo 不支持写
	return nil
}

func NewMMgo(param string) (*mmgo, error) {
	r := new(mmgo)
	e := r.init(param)
	if e != nil {
		return nil, e
	}
	return r, nil
}
