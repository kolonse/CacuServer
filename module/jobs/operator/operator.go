package operator

import (
	"fmt"
	"strings"

	"github.com/kolonse/CacuServer/module/jobs/operator/db"
)

type DB interface {
	Count() (int, error)
	Reset()
	Read() (interface{}, error)
	Write() error
}

func ParseDB(str string) (DB, error) {
	index := strings.Index(str, " ")
	if index == -1 {
		return nil, fmt.Errorf("配置格式不正确")
	}
	rtype := str[0:index]
	switch rtype {
	case "mysql":
		return db.NewMMysql(str[index+1:])
	case "mongo":
		return db.NewMMgo(str[index+1:])
	}
	return nil, fmt.Errorf("[%v] unkown db type [%v]", str, rtype)
}
