package module

import (
	"fmt"
	"testing"
	"time"

	"github.com/kolonse/CacuServer/conf"
)

type Log struct {
}

func (l *Log) Info(str string, args ...interface{}) {
	fmt.Printf(str, args...)
	fmt.Println()
}

func (l *Log) Debug(str string, args ...interface{}) {
	fmt.Printf(str, args...)
	fmt.Println()
}

func (l *Log) Error(str string, args ...interface{}) {
	fmt.Printf(str, args...)
	fmt.Println()
}

func TestLoadCfg(t *testing.T) {
	conf.Cfg.ParseFile("../conf/conf.kcfg")
	conf.Cfg.Child("JobPath").SetString("../conf/job.kcfg")
	Module.SetLogger(new(Log))
	Module.Init()
	sig := make(chan bool, 1)
	go func() {
		Module.Run()
		sig <- true
	}()
	go func() {
		time.Sleep(10 * time.Second)
		Module.Stop()
	}()
	<-sig
}
