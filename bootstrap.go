package bone

//Bootstrap 接口
type Bootstrap interface {
	Run() //消息循环的空闲例程
}

//NewBootstrap 返回1个Bootstrap
func NewBootstrap(template string, logger BLogger, provider []interface{}) Bootstrap {
	if logger == nil {
		logger = &emptyBLogger{}
	}
	if template == "" {
		logger.Error("empty template")
		return nil
	}
	doc := xmlDoc{}
	if err := doc.Parse(template); err != nil {
		logger.Error(err.Error())
		return nil
	}
	td := newTaskDispatcher()
	wm := newWatcherManager()
	defaultProvder := make([]interface{}, 3)
	defaultProvder[0] = TaskDispatcher(td)
	defaultProvder[1] = WatcherManager(wm)
	defaultProvder[2] = logger

	ij := &Injector{
		Parent:   nil,
		Provider: append(append([]interface{}(nil), provider...), defaultProvder...),
	}
	bt := &bootstrap{
		Bone:           NewBone(ij, nil, doc),
		topLevel:       NewView("View"),
		watcherManager: wm,
		dispatcher:     td,
	}
	//调用load来解析
	bt.watcherManager.SetMode(false) //load时设置为同步模式
	if err := bt.load(bt.topLevel); err != nil {
		logger.Error(err.Error())
		return nil
	}
	bt.watcherManager.SetMode(true) //正常显示时为异步模式 在之后支持for if时也会临时设置为同步模式
	return bt
}

type bootstrap struct {
	*Bone
	topLevel       *View
	watcherManager *watcherManager
	dispatcher     *taskDispatcher
}

func (bs *bootstrap) Run() {
	bs.dispatcher.Run()
	bs.watcherManager.Run()
}
