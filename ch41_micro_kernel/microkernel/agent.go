package microkernel

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

//常量 iota
const (
	Waiting = iota
	Running
)

//错误
var WrongStateError = errors.New("can not take the operation in the current state")

//收集器错误
type CollectorsError struct {
	CollectorErrors []error
}

//收集器错误处理
func (ce CollectorsError) Error() string {
	var strs []string
	for _, err := range ce.CollectorErrors {
		strs = append(strs, err.Error())
	}
	return strings.Join(strs, ";")
}

//-----------------------------------------------------

type Event struct {
	Source  string
	Content string
}

type EventReceiver interface {
	OnEvent(evt Event)
}

//微内核架构模式

//收集信息接口
type Collector interface {
	// 收集的信息 回传给的实现了 OnEvent 接口的对象
	Init(evtReceiver EventReceiver) error
	//上下文 为了cancel掉不同协程的任务
	Start(agtCtx context.Context) error
	Stop() error
	Destory() error
}

// agent 对象 管理收集器的
type Agent struct {
	//收集器 map
	collectors map[string]Collector
	//chan 缓冲区
	evtBuf chan Event

	//取消上下文
	cancel context.CancelFunc
	//上下文 子节点
	ctx context.Context

	//状态
	state int
}

//单独协程 收集器启动会上报信息  收集够10条打印一次
func (agt *Agent) EventProcessGroutine() {
	var evtSeg [10]Event
	for {
		for i := 0; i < 10; i++ {
			select {
			//从chan接收数据
			case evtSeg[i] = <-agt.evtBuf:
			case <-agt.ctx.Done():
				return
			}
		}
		fmt.Println(evtSeg)
	}

}

//new agent 管理收集器的
func NewAgent(sizeEvtBuf int) *Agent {
	agt := Agent{
		collectors: map[string]Collector{},
		evtBuf:     make(chan Event, sizeEvtBuf),
		state:      Waiting, //等待
	}

	return &agt
}

//注册收集器
func (agt *Agent) RegisterCollector(name string, collector Collector) error {
	//等待状态 才能注册
	if agt.state != Waiting {
		return WrongStateError
	}
	agt.collectors[name] = collector
	//注册后会调用对应单个收集器初始化
	return collector.Init(agt)
}

func (agt *Agent) startCollectors() error {
	var err error
	var errs CollectorsError //收集器错误
	var mutex sync.Mutex

	//启动 所有注册的收集器
	for name, collector := range agt.collectors {
		go func(name string, collector Collector, ctx context.Context) {
			defer func() {
				mutex.Unlock()
			}()
			err = collector.Start(ctx) //是否成功启动
			mutex.Lock()
			if err != nil {
				errs.CollectorErrors = append(errs.CollectorErrors,
					errors.New(name+":"+err.Error()))
			}
		}(name, collector, agt.ctx)
	}
	if len(errs.CollectorErrors) == 0 {
		return nil
	}
	return errs
}

func (agt *Agent) stopCollectors() error {
	var err error
	var errs CollectorsError
	//停止agent所有收集器
	for name, collector := range agt.collectors {
		if err = collector.Stop(); err != nil {
			errs.CollectorErrors = append(errs.CollectorErrors,
				errors.New(name+":"+err.Error()))
		}
	}
	if len(errs.CollectorErrors) == 0 {
		return nil
	}

	return errs
}

func (agt *Agent) destoryCollectors() error {
	var err error
	var errs CollectorsError
	//注销agent所有收集器
	for name, collector := range agt.collectors {
		if err = collector.Destory(); err != nil {
			errs.CollectorErrors = append(errs.CollectorErrors,
				errors.New(name+":"+err.Error()))
		}
	}
	if len(errs.CollectorErrors) == 0 {
		return nil
	}
	return errs
}

func (agt *Agent) Start() error {
	if agt.state != Waiting {
		return WrongStateError
	}
	agt.state = Running
	//agt.ctx上下文
	agt.ctx, agt.cancel = context.WithCancel(context.Background())

	//起一个协程 收集器启动会上报信息  收集够10条打印一次
	go agt.EventProcessGroutine()

	//启动所有收集器
	return agt.startCollectors()
}

func (agt *Agent) Stop() error {
	if agt.state != Running {
		return WrongStateError
	}
	agt.state = Waiting
	//context.cancel()取消上下文
	agt.cancel()

	return agt.stopCollectors()
}

func (agt *Agent) Destory() error {
	if agt.state != Waiting {
		return WrongStateError
	}

	return agt.destoryCollectors()
}

func (agt *Agent) OnEvent(evt Event) {
	//接收收集的信息 并且放入缓冲区chan
	agt.evtBuf <- evt
}
