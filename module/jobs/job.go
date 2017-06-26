package jobs

import (
	kcfg "github.com/kolonse/kolonsecfg"
)

type logger interface {
	Error(string, ...interface{})
}

type Job interface {
	// 初始化任务
	Init(node *kcfg.Node) error
	// 检测是否达到统计时间
	OnTime() <-chan bool
	// 开始计算
	Cacu() error
	// 设置日志接口
	SetLogger(logger)
}

func NewJob() *cacujob {
	return new(cacujob)
}
