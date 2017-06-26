// CacuServer project main.go
package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"

	"github.com/kolonse/CacuServer/app"
	"github.com/kolonse/CacuServer/conf"
	"github.com/kolonse/CacuServer/log"
	"github.com/kolonse/CacuServer/module"
	"github.com/kolonse/CacuServer/rpc"
)

var cfgPath = flag.String("cfgPath", "./conf/conf.kcfg", "-cfgPath=<string> configure file path,default ./conf/conf.kcfg")

func main() {
	flag.Parse()
	// 加载配置
	conf.Cfg.ParseFile(*cfgPath)
	conf.MemoryLimit = conf.Cfg.Child("MemoryLimit").GetInt()
	conf.OnceSleepTime = int(conf.Cfg.Child("OnceSleepTime").GetInt())
	conf.ReadCountLimit = int(conf.Cfg.Child("ReadCountLimit").GetInt())
	// 初始化 Log
	log.SetLogPath(conf.Cfg.Child("Log.Path").GetString(),
		conf.Cfg.Child("Log.Name").GetString())
	log.Run()
	log.Logger.Info("日志组件加载完成")
	defer log.Logger.Close()
	rpc.DefaultRPC.SetLogger(log.Logger)
	module.Module.SetLogger(log.Logger)
	app.Init(rpc.DefaultRPC, module.Module)
	log.SetLevel(conf.Cfg.Child("Log.Level").GetString())
	go func() {
		http.HandleFunc("/goroutines", func(w http.ResponseWriter, r *http.Request) {
			num := strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
			w.Write([]byte(num))
		})
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()
	app.Go()
}
