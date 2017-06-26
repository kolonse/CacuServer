package rpc

import (
	"fmt"
	"reflect"
)

const (
	// MaxCallLimit 最大调用限制,如果超过就会有阻塞
	MaxCallLimit = 10000
)

var (
	seq = 0
)

// KRPCError 错误描述
type KRPCError struct {
	err string
}

func (err *KRPCError) Error() string {
	return err.err
}

type logger interface {
	Info(string, ...interface{})
}

type retInfo struct {
	key  int
	cb   interface{}
	err  error
	args []interface{}
}

type callInfo struct {
	key  int
	f    interface{}
	cb   interface{}
	args []interface{}
	ret  chan retInfo
}

// RPC 结构,用户处理内部RPC任务
type RPC struct {
	exit chan bool
	q    chan callInfo
	r    chan retInfo
	k2v  map[interface{}]int
	v2k  map[int]interface{}
	k2f  map[int]callInfo
	log  logger
}

// Init 初始化接口
func (r *RPC) Init() {
	r.exit = make(chan bool, 1)
	r.q = make(chan callInfo, MaxCallLimit)
	r.r = make(chan retInfo, MaxCallLimit)
	r.k2v = make(map[interface{}]int)
	r.v2k = make(map[int]interface{})
	r.k2f = make(map[int]callInfo)
	r.log.Info("RPC Init 完成")
}

// SetLogger 设置日志实例接口,接口必须实现 Info,Debug,Warn,Error接口
func (r *RPC) SetLogger(log logger) {
	r.log = log
}

// Run 运行函数
func (r *RPC) Run() error {
	r.log.Info("RPC 服务启动完成")
	sig1 := make(chan int, 1)
	sig2 := make(chan int, 1)
	sig := make(chan int, 1)
	go func() {
		<-r.exit
		sig1 <- 1
		sig2 <- 2
	}()
	go func() {
		defer func() {
			sig <- 1
		}()
		for {
			select {
			case <-sig1:
				return
			case cf, ok := <-r.q:
				if !ok {
					return
				}
				// 处理远程调用
				cf.args = append(cf.args, func(err error, args ...interface{}) {
					ret := retInfo{
						key:  cf.key,
						err:  err,
						args: args,
						cb:   cf.cb,
					}
					r.r <- ret
				})
				var in []reflect.Value
				var v reflect.Value
				for j := 0; j < len(cf.args); j++ {
					if cf.args[j] == nil {
						v = reflect.Zero(reflect.TypeOf(cf.args[j]))
					} else {
						v = reflect.ValueOf(cf.args[j])
					}
					in = append(in, v)
				}
				reflect.ValueOf(cf.f).Call(in)
			}
		}
	}()
	go func() {
		defer func() {
			sig <- 2
		}()
		for {
			select {
			case <-sig2:
				return
			case ret, ok := <-r.r:
				if !ok {
					return
				}
				// 处理回调通知
				var in []reflect.Value
				var v reflect.Value
				if ret.err == nil {
					v = reflect.ValueOf(&KRPCError{})
				} else {
					v = reflect.ValueOf(ret.err)
				}
				in = append(in, v)
				for j := 0; j < len(ret.args); j++ {
					if ret.args[j] == nil {
						v = reflect.Zero(reflect.TypeOf(ret.args[j]))
					} else {
						v = reflect.ValueOf(ret.args[j])
					}
					in = append(in, v)
				}
				reflect.ValueOf(ret.cb).Call(in)
			}
		}
	}()

	<-sig
	<-sig
	return nil
}

// Stop 停止rpc服务 运行
func (r *RPC) Stop() {
	r.exit <- true
}

// Exit RPC退出处理函数
func (r *RPC) Exit(err error) {
	r.log.Info("RPC 服务退出完成")
	return
}

func (r *RPC) info(fmt string, args ...interface{}) {
	if r.log != nil {
		r.log.Info(fmt, args...)
	}
	return
}

// Call rpc调用接口,需要传入回调接口,接口格式 func (error,args ...interface{})
func (r *RPC) Call(key interface{}, args ...interface{}) {
	i, exist := r.k2v[key]
	if !exist {
		panic(fmt.Errorf("%v method not register", key))
	}
	var cf callInfo
	cf = r.k2f[i]
	cf.args = args
	cf.ret = make(chan retInfo, 1)
	var cb interface{}
	if len(args) != 0 {
		cb = args[len(args)-1]
		switch reflect.TypeOf(cb).Kind() {
		case reflect.Func:
			cf.args = cf.args[0 : len(cf.args)-1]
		default:
			cb = func(error, ...interface{}) {}
		}
	}
	cf.cb = cb
	r.q <- cf
}

// Notify rpc通知类型接口,不需要传入回调
func (r *RPC) Notify(key interface{}, args ...interface{}) {
	panic(fmt.Errorf("not support"))
}

// Register rpc注册接口
func (r *RPC) Register(key interface{}, f interface{}) {
	r.k2v[key] = seq
	r.v2k[seq] = key

	cf := callInfo{
		key: seq,
		f:   f,
	}
	r.k2f[seq] = cf
	seq = seq + 1
	return
}

// DefaultRPC 默认rpc接口实例
var DefaultRPC = new(RPC)
