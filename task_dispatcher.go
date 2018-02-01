package bone

import (
	"container/list"
	"sync"
)

//TaskDispatcher 保证在UI线程执行
type TaskDispatcher interface {
	Invoke(do func())
}

//newTaskDispatcher 返回1个dispatcher
func newTaskDispatcher() *taskDispatcher {
	return &taskDispatcher{
		Locker: &sync.Mutex{},
		call:   list.New(),
	}
}

type taskDispatcher struct {
	sync.Locker
	call *list.List
}

func (td *taskDispatcher) Invoke(do func()) {
	td.Lock()
	defer td.Unlock()
	td.call.PushBack(do)
}

func (td *taskDispatcher) Peek() func() {
	td.Lock()
	defer td.Unlock()
	if td.call.Len() > 0 {
		front := td.call.Front()
		ret := front.Value.(func())
		td.call.Remove(front)
		return ret
	}
	return nil

}

func (td *taskDispatcher) Run() {
	for {
		do := td.Peek()
		if do == nil {
			break
		}
		do()
	}
}
