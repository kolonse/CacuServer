package conf

import (
	conf "github.com/kolonse/kolonsecfg"
)

var (
	Cfg                  = conf.NewCfg()      // 配置
	MemoryLimit    int64 = 1024 * 1024 * 1024 //default 1G
	OnceSleepTime        = 10 * 1000          // default 10s
	ReadCountLimit       = 10000
)
