package module

import (
	"time"

	"github.com/kolonse/CacuServer/conf"
	"github.com/kolonse/CacuServer/module/jobs"
	kcfg "github.com/kolonse/kolonsecfg"
)

var Module = new(StatisticsModule)

type logger interface {
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Error(string, ...interface{})
}

type StatisticsModule struct {
	jobPath string
	jobCfg  *kcfg.Cfg
	log     logger
	jobs    []jobs.Job
	cj      chan int
	exit    chan bool
}

// SetLogger 设置日志实例接口,接口必须实现 Info,Debug,Warn,Error接口
func (m *StatisticsModule) SetLogger(log logger) {
	m.log = log
}

func (m *StatisticsModule) Init() {
	m.cj = make(chan int, 1)
	m.exit = make(chan bool, 1)
	m.jobPath = conf.Cfg.Child("JobPath").GetString()
	m.jobCfg = kcfg.NewCfg()
	m.loadCfg()
	m.loadJobs()
	m.log.Info("StatisticsModule Init 完成")
}

func (m *StatisticsModule) Run() error {
	m.log.Info("StatisticsModule 服务启动完成")
	for {
		select {
		case <-m.exit:
			return nil
		case index, ok := <-m.cj:
			if !ok {
				return nil
			}
			m.process(index)
		}
	}
	return nil
}

func (m *StatisticsModule) Stop() {
	m.exit <- true
}

func (m *StatisticsModule) Exit(error) {
	m.log.Info("服务退出完成")
}

func (m *StatisticsModule) loadCfg() {
	m.jobCfg.ParseFile(m.jobPath)
}

func (m *StatisticsModule) loadJobs() {
	js := m.jobCfg.Childs("Jobs")
	for i := 0; i < len(js); i++ {
		job := jobs.NewJob()
		job.SetLogger(m.log)
		err := job.Init(js[i])
		if err != nil {
			m.log.Error("Jobs %v 加载未成功,err:%v", i, err.Error())
			continue
		}
		go func(index int, j jobs.Job) {
			for {
				select {
				case <-j.OnTime():
					// 通知模块有任务需要进行统计
					m.cj <- index
				}
			}
		}(i, job)
		m.jobs = append(m.jobs, job)
	}
}

/**
*	process 处理步骤
*	1. 系统剩余内存 = 当前系统内存 - 当前占用内存 < 需求内存 - 进程内存
*	2. 检测 CPU 状态, CPU < 设定最大值
*	3. 获取数据 Count
*	4. 分批获取数据,每次的数据不能超过 Count 设定值
*	5. 计算好后,进行存储
 */
func (m *StatisticsModule) process(index int) {
	m.log.Info("开始处理任务 %v", index)
	b := time.Now().UnixNano()
	job := m.jobs[index]
	err := job.Cacu()
	if err != nil {
		m.log.Error("任务 %v 处理异常,err:%v", index, err.Error())
		return
	}
	m.log.Info("处理任务 %v 结束,花费时间:%v", index, (time.Now().UnixNano()-b)/1000000.0)
	time.Sleep(time.Duration(conf.OnceSleepTime) * time.Second)
}
