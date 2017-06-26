package jobs

import (
	"github.com/kolonse/CacuServer/module/jobs/operator"
	"github.com/kolonse/CacuServer/script"
	kcfg "github.com/kolonse/kolonsecfg"
	"github.com/robfig/cron"
)

type cacujob struct {
	timeStr     string
	countStr    string
	readStr     string
	gobjectStr  []string
	cacuStr     []string
	storeStr    []string
	chanOnTime  chan bool
	countOpt    operator.DB
	readOpt     operator.DB
	gobjectOpt  []operator.DB
	cacuScript  []script.Script
	storeScript []script.Script
	log         logger
}

func (j *cacujob) Init(node *kcfg.Node) error {
	j.timeStr = node.Child("Time").GetString()
	j.countStr = node.Child("Count").GetString()
	j.readStr = node.Child("Array").GetString()
	ns := node.Childs("GObject")
	for i := 0; i < len(ns); i++ {
		j.gobjectStr = append(j.gobjectStr, ns[i].GetString())
	}
	ns = node.Childs("Cacu")
	for i := 0; i < len(ns); i++ {
		j.cacuStr = append(j.cacuStr, ns[i].GetString())
		spt := script.NewRuntimeScript()
		spt.Parse(j.cacuStr[i])
		j.cacuScript = append(j.cacuScript, spt)
	}
	ns = node.Childs("Store")
	for i := 0; i < len(ns); i++ {
		j.storeStr = append(j.storeStr, ns[i].GetString())
		spt := script.NewRuntimeScript()
		spt.Parse(j.storeStr[i])
		j.storeScript = append(j.storeScript, spt)
	}
	countOpt, err := operator.ParseDB(j.countStr)
	if err != nil {
		return err
	}
	j.countOpt = countOpt

	j.readOpt, err = operator.ParseDB(j.readStr)
	if err != nil {
		return err
	}

	for _, v := range j.gobjectStr {
		opt, e := operator.ParseDB(v)
		if e != nil {
			return e
		}
		j.gobjectOpt = append(j.gobjectOpt, opt)
	}
	j.chanOnTime = make(chan bool, 1)

	c := cron.New()
	c.AddFunc(j.timeStr, func() {
		j.chanOnTime <- true
	})
	c.Start()
	return nil
}

func (j *cacujob) SetLogger(log logger) {
	j.log = log
}

func (j *cacujob) OnTime() <-chan bool {
	return j.chanOnTime
}

func (j *cacujob) Cacu() error {
	count, e := j.countOpt.Count()
	if e != nil {
		return e
	}
	rcount := 0
	if count == 0 {
		return nil
	}
	// 读取所有的  GObject 数据
	datamap := make(map[string]interface{})
	gobjects, e := j.readGObject()
	if e != nil {
		return e
	}
	datamap["GObject"] = make([]interface{}, len(gobjects))
	for i, v := range gobjects {
		datamap["GObject"].([]interface{})[i] = v
	}
	j.readOpt.Reset()
	datamap["Array"] = make([]interface{}, 1)
	for rcount < count {
		record, err := j.readOpt.Read()
		if err != nil {
			return err
		}
		if record == nil {
			break
		}
		datamap["Array"].([]interface{})[0] = record.([]interface{})
		// 这里进行 Cacu 计算
		/**
		*	Cacu 的计算都是基于record循环遍历进行的,所以这里必须要有个 record 循环,然后调用 Cacu 方法
		 */
		for i, _ := range record.([]interface{}) {
			for t, v := range j.cacuScript {
				v.SetCallArgs(datamap, i)
				ri, ei := v.Call()
				if ei != nil {
					j.log.Error("计算脚本 %v,数据 %v 错误,忽略该数据,%v", t, record.([]interface{})[i], ei.Error())
					continue
				}
				datamap = ri.(map[string]interface{})
			}
		}
		rcount = rcount + len(record.([]interface{}))
		if len(record.([]interface{})) == 0 {
			break
		}
	}
	//调用存储接口
	for i, v := range j.storeScript {
		v.SetCallArgs(datamap)
		f, err := v.Call()
		if err != nil {
			j.log.Error("数据存储 %v 错误,忽略该存储,%v", i, err.Error())
			continue
		}
		f.(func(func(string)))(func(str string) {
			//创建存储实例
			opt, e := operator.ParseDB(str)
			if e != nil {
				j.log.Error("数据存储 %v 错误,忽略该存储,%v", str, e.Error())
				return
			}
			e = opt.Write()
			if e != nil {
				j.log.Error("数据存储 %v 错误,忽略该存储,%v", str, e.Error())
				return
			}
		})
	}
	return nil
}

func (j *cacujob) readGObject() ([]map[string]interface{}, error) {
	ret := make([]map[string]interface{}, 0)
	for _, v := range j.gobjectOpt {
		v.Reset()
		r, e := j.readAll(v)
		if e != nil {
			return nil, e
		}
		if r == nil {
			ret = append(ret, make(map[string]interface{}))
			continue
		}
		ret = append(ret, r.(map[string]interface{}))
	}
	return ret, nil
}

func (j *cacujob) readAll(opt operator.DB) (interface{}, error) {
	var ret interface{}
	for {
		r, e := opt.Read()
		if e != nil {
			return nil, e
		}
		count := 0
		if r == nil {
			return ret, nil
		}
		switch r.(type) {
		case []interface{}:
			if ret == nil {
				ret = make([]interface{}, 0)
			}
			count = len(r.([]interface{}))
			if count == 0 {
				return ret, nil
			}

			ret = append(ret.([]interface{}), r.([]interface{})...)
		case map[string]interface{}:
			if ret == nil {
				ret = make(map[string]interface{})
			}
			count = len(r.(map[string]interface{}))
			if count == 0 {
				return ret, nil
			}
			for k, v := range r.(map[string]interface{}) {
				ret.(map[string]interface{})[k] = v
			}
		}
	}
	return ret, nil
}
