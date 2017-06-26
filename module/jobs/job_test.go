package jobs

import (
	"fmt"
	"testing"
	"time"

	kcfg "github.com/kolonse/kolonsecfg"
)

type Log struct {
}

func (log *Log) Error(str string, args ...interface{}) {
	fmt.Printf(str+"\n", args...)
}

func TestOnTime(t *testing.T) {
	j := new(cacujob)
	j.Init(kcfg.NewCfg().ParseFile("../../conf/job.kcfg").Child("Jobs"))
	l := new(Log)
	j.SetLogger(l)
	go func() {
		for {
			select {
			case <-j.OnTime():
				fmt.Println("ontime")
				e := j.Cacu()
				fmt.Println(e)
				return
			}
		}
	}()
	time.Sleep(5 * time.Second)
}
