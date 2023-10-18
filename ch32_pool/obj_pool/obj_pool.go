package object_pool

import (
	"errors"
	"time"
)

//池子里的对象
type ReusableObj struct {
}

type ObjPool struct {
	bufChan chan *ReusableObj //用于缓冲可重用对象 chan存放ReusableObj
}

//回对象池 参数长度
func NewObjPool(numOfObj int) *ObjPool {
	objPool := ObjPool{}
	//有缓存的
	objPool.bufChan = make(chan *ReusableObj, numOfObj)
	for i := 0; i < numOfObj; i++ {
		objPool.bufChan <- &ReusableObj{}
	}
	return &objPool
}

//ObjPool对象指针的方法 返回对象指针
func (p *ObjPool) GetObj(timeout time.Duration) (*ReusableObj, error) {
	select {
	case ret := <-p.bufChan:
		return ret, nil
	case <-time.After(timeout): //超时控制
		return nil, errors.New("time out")
	}

}

// 新增 对象到池子里
func (p *ObjPool) ReleaseObj(obj *ReusableObj) error {
	select {
	case p.bufChan <- obj:
		return nil
	default:
		//放不进去 对象池 chan长度限制
		return errors.New("overflow")
	}
}
