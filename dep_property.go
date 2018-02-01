package bone

import (
	"sort"
)

//WatchedTarget 可以被Watcher监听的接口
type WatchedTarget interface {
	Get() interface{}
	SetFrom(WatchedTarget, interface{})
	Equal(interface{}) bool
	AddSub(watcher *Watcher)
}

//Watcher 观测属性的变化
type Watcher struct {
	source WatchedTarget
	target WatchedTarget
	//update       func(sender interface{}, source interface{}, old, new interface{})
	computer func(sender interface{}, new interface{}) interface{} //用于支持计算属性
	queue    *watcherQueue
}

//Update 需要更新时调用update函数
func (wr *Watcher) Update() {
	newValue := wr.source.Get()
	if wr.computer != nil {
		newValue = wr.computer(wr, newValue)
	}
	if wr.target.Equal(newValue) {
		return
	}
	wr.target.SetFrom(wr.source, newValue)
}

//Run 异步执行
func (wr *Watcher) Run() {
	if wr.queue != nil {
		wr.queue.Push(wr)
	} else {
		wr.Update()
	}
}

//Dependence 依赖
type Dependence struct {
	subs []*Watcher
}

//Notify 通知
func (dp *Dependence) Notify(source WatchedTarget) {
	for _, v := range dp.subs {
		if v.target != source {
			v.Run()
		}
	}
}

//AddSub 添加订阅 当支持if 和 for指令时 可能需要实现RemoveSub
func (dp *Dependence) AddSub(watcher *Watcher) {
	if watcher == nil {
		return
	}
	dp.subs = append(dp.subs, watcher)
}

//DepProperty 用于数据绑定
type DepProperty interface {
	Set(value interface{})
	Get() interface{}
}

type depProperty struct {
	Dependence
	value  interface{}
	setted func(sender interface{}, value interface{})
	getted func(sender interface{}, value interface{})
}

//NewDepProperty 返回1个依赖对象
func NewDepProperty(value interface{}, setted, getted func(sender interface{}, value interface{})) DepProperty {
	return &depProperty{
		value:  value,
		setted: setted,
		getted: getted,
	}
}

//Set 设置属性值
func (dp *depProperty) Set(value interface{}) {
	dp.SetFrom(nil, value)
}

//Get 获取属性值
func (dp *depProperty) Get() interface{} {
	if dp.getted != nil {
		dp.getted(dp, dp.value)
	}
	return dp.value
}

func (dp *depProperty) SetFrom(source WatchedTarget, value interface{}) {
	if dp.Equal(value) {
		return
	}
	dp.value = value
	if dp.setted != nil {
		dp.setted(dp, dp.value)
	}
	dp.Notify(source)
}

//Equal value是否相等
func (dp *depProperty) Equal(value interface{}) bool {
	return dp.value == value
}

//WatcherManager 用于管理Watcher
type WatcherManager interface {
	Create(source, target WatchedTarget,
		computer func(sender interface{}, new interface{}) interface{}) *Watcher
	SetMode(async bool)
	Mode() bool
}

func newWatcherManager() *watcherManager {
	return &watcherManager{
		watcherQueue: watcherQueue{
			async: true,
			slice: nil,
			exist: make(map[*Watcher]bool),
		},
	}
}

type watcherManager struct {
	watcherQueue
}

func (wm *watcherManager) Create(source, target WatchedTarget,
	computer func(sender interface{}, new interface{}) interface{}) *Watcher {
	wr := &Watcher{
		source:   source,
		target:   target,
		computer: computer,
		queue:    &wm.watcherQueue,
	}
	return wr
}

type watcherQueue struct {
	async bool //true 则入队  false直接执行
	slice []*Watcher
	exist map[*Watcher]bool
}

func (wq *watcherQueue) SetMode(async bool) {
	if wq.async == async {
		return
	}

	for len(wq.slice) > 0 {
		wq.Run()
	}
	wq.async = async
}

func (wq *watcherQueue) Mode() bool {
	return wq.async
}

//Push Watch入队 自动去重
func (wq *watcherQueue) Push(watcher *Watcher) {
	if !wq.async {
		watcher.Update()
		return
	}
	if _, ok := wq.exist[watcher]; ok {
		return
	}
	wq.exist[watcher] = true
	wq.slice = append(wq.slice, watcher)
}

//Run 执行异步watch
func (wq *watcherQueue) Run() {
	for _, v := range wq.slice {
		v.Update()
	}

	wq.slice = wq.slice[:0]
	for k := range wq.exist {
		delete(wq.exist, k)
	}
}

//SliceDepProperty 用于数据绑定
type SliceDepProperty interface {
	SetSlice(value []interface{})
	GetSlice() []interface{}
	Push(item interface{}) int
	Pop() interface{}
	Unshift(item interface{}) int
	Shift() interface{}
	Splice(start, howmany int, items ...interface{}) []interface{}
}

//SliceDepProperty 数组依赖属性
type sliceDepProperty struct {
	Dependence
	value  []interface{}
	setted func(sender interface{}, value []interface{})
	getted func(sender interface{}, value []interface{})
}

//SliceDepProperty 接口
func (dp *sliceDepProperty) SetSlice(value []interface{}) {
	dp.SetFrom(nil, value)
}

func (dp *sliceDepProperty) GetSlice() []interface{} {
	return dp.Get().([]interface{})
}

func (dp *sliceDepProperty) Push(item interface{}) int {
	dp.value = append(dp.value, item)
	dp.notifyIfSet(nil)
	return len(dp.value)
}

func (dp *sliceDepProperty) Pop() interface{} {
	count := len(dp.value)
	if count <= 0 {
		return nil
	}
	tail := dp.value[count-1]
	dp.value = dp.value[0 : count-1]
	dp.notifyIfSet(nil)
	return tail
}

func (dp *sliceDepProperty) Unshift(item interface{}) int {
	value := make([]interface{}, len(dp.value)+1)
	value[0] = item
	for i := 0; i < len(dp.value); i++ {
		value[i+1] = dp.value[i]
	}
	dp.value = value
	dp.notifyIfSet(nil)
	return len(dp.value)
}

func (dp *sliceDepProperty) Shift() interface{} {
	count := len(dp.value)
	if count <= 0 {
		return nil
	}
	head := dp.value[0]
	dp.value = dp.value[1:]
	dp.notifyIfSet(nil)
	return head
}

//用于调用排序方法
type sortInterface struct {
	value  []interface{}
	method func(i, j interface{}) bool
}

func (si sortInterface) Len() int {
	return len(si.value)
}

func (si sortInterface) Less(i, j int) bool {
	return si.method(si.value[i], si.value[j])
}

func (si sortInterface) Swap(i, j int) {
	si.value[i], si.value[j] = si.value[j], si.value[i]
}

func (dp *sliceDepProperty) Sort(method func(i, j interface{}) bool) {
	if len(dp.value) <= 0 {
		return
	}
	sort.Sort(sortInterface{value: dp.value, method: method})
	dp.notifyIfSet(nil)
}

func (dp *sliceDepProperty) Splice(start, howmany int, items ...interface{}) []interface{} {
	count := len(dp.value)
	if start < 0 || howmany < 0 || start+howmany > count {
		return nil
	}
	value := make([]interface{}, count-howmany+len(items))
	j := 0
	for i := 0; i < start; i++ {
		value[j] = dp.value[i]
		j++
	}

	for i := 0; i < len(items); i++ {
		value[j] = items[i]
		j++
	}

	for i := start + howmany; i < count; i++ {
		value[j] = dp.value[i]
		j++
	}
	removed := dp.value[start : start+howmany]
	dp.value = value
	dp.notifyIfSet(nil)
	return removed
}

func (dp *sliceDepProperty) Reverse() {
	count := len(dp.value)
	if count <= 1 {
		return
	}
	for i, j := 0, count-1; i < j; i, j = i+1, j-1 {
		dp.value[i], dp.value[j] = dp.value[j], dp.value[i]
	}
	dp.notifyIfSet(nil)
}

func (dp *sliceDepProperty) notifyIfSet(source WatchedTarget) {
	if dp.setted != nil {
		dp.setted(dp, dp.value)
	}
	dp.Notify(source)
}

//Equal 数组是否相同
func (dp *sliceDepProperty) Equal(value interface{}) bool {
	left := dp.value
	right := value.([]interface{})

	if len(left) != len(right) {
		return false
	}
	count := len(left)
	for i := 0; i < count; i++ {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

//SetFrom 设置属性值
func (dp *sliceDepProperty) SetFrom(source WatchedTarget, value interface{}) {
	if dp.Equal(value) {
		return
	}
	dp.value = append([]interface{}(nil), value)
	dp.notifyIfSet(source)
}

//Get 获取属性值
func (dp *sliceDepProperty) Get() interface{} {
	if dp.getted != nil {
		dp.getted(dp, dp.value)
	}
	return dp.value
}
