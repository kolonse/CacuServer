package rpc

import (
	"fmt"
	"testing"
	"time"
)

type Log struct {
}

func (l *Log) Info(str string, args ...interface{}) {
	fmt.Printf(str, args...)
}

func TestCall(t *testing.T) {
	rpc := &RPC{}
	rpc.Init()
	rpc.Register("f1", func(n1 int64, cb func(error, ...interface{})) {
		cb(nil, time.Now().UnixNano()-n1)
	})
	sig := make(chan bool, 1)
	go func() {
		rpc.Run()
		sig <- true
	}()
	count := 0
	times := 500000
	timecost := int64(0)
	for i := 0; i < times; i++ {
		rpc.Call("f1", time.Now().UnixNano(), func(err error, data int64) {
			count = count + 1
			timecost = timecost + data
			if count >= times {
				t.Log("Count:", count, "time:", timecost, "avg:", timecost/int64(times), "ns/op")
				rpc.Stop()
			}
		})
	}
	<-sig
}

func BenchmarkCall(b *testing.B) {
	rpc := &RPC{}
	rpc.Init()
	rpc.Register("f1", func(n1 int64, cb func(error, ...interface{})) {
		cb(nil, time.Now().UnixNano()-n1)
	})
	sig := make(chan bool, 1)
	go func() {
		rpc.Run()
		sig <- true
	}()
	count := 0
	times := b.N
	for i := 0; i < times; i++ {
		rpc.Call("f1", time.Now().UnixNano(), func(err error, data int64) {
			count = count + 1
			if count >= times {
				rpc.Stop()
			}
		})
	}
	<-sig
}
