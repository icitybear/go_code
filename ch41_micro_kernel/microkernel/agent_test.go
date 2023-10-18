package microkernel

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

//demo收集器
type DemoCollector struct {
	evtReceiver EventReceiver

	agtCtx   context.Context
	stopChan chan struct{}

	name    string
	content string
}

func NewCollect(name string, content string) *DemoCollector {
	return &DemoCollector{
		stopChan: make(chan struct{}),
		name:     name,
		content:  content,
	}
}

//DemoCollector 的init
func (c *DemoCollector) Init(evtReceiver EventReceiver) error {
	fmt.Println("agent RegisterCollector after initialize collector", c.name)
	c.evtReceiver = evtReceiver
	return nil
}

func (c *DemoCollector) Start(agtCtx context.Context) error {
	fmt.Println("start DemoCollector", c.name)
	for {
		select {
		//收到收集器的上下文的cancel
		case <-agtCtx.Done():
			//给DemoCollector stopChan的缓冲区放消息
			fmt.Println("DemoCollector add stopChan", c.name)
			c.stopChan <- struct{}{}
			break
		default:
			//每隔一段时间 上报收集信息
			time.Sleep(time.Millisecond * 50)
			//关于OnEvent的实现 只有agent
			c.evtReceiver.OnEvent(Event{c.name, c.content})
		}
	}
}

func (c *DemoCollector) Stop() error {
	fmt.Println("stop DemoCollector", c.name)
	select {
	//收到 stopChan的缓冲区信息
	case <-c.stopChan:
		return nil
	case <-time.After(time.Second * 1):
		return errors.New("failed to stop for timeout")
	}
}

func (c *DemoCollector) Destory() error {
	fmt.Println(c.name, "released resources.")
	return nil
}

func TestAgent(t *testing.T) {
	agt := NewAgent(100)        //100是agent的chan缓冲区大小
	c1 := NewCollect("c1", "1") //初始化收集器 name content stopChan
	c2 := NewCollect("c2", "2")

	agt.RegisterCollector("c1", c1) //注册到agent收集器容器里  注册后会调用对应单个收集器初始化Init
	agt.RegisterCollector("c2", c2)
	//开启
	if err := agt.Start(); err != nil {
		fmt.Printf("start error %v\n", err)
	}
	fmt.Println(agt.Start()) //agent状态state已经是开启的了
	time.Sleep(time.Second * 1)

	//agent的停止与注销
	agt.Stop() //agt.cancel() 并 循环stop所有的收集器
	agt.Destory()
}
