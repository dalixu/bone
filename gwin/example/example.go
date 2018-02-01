package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/dalixu/bone"
	"github.com/dalixu/bone/gwin"
	"github.com/lxn/win"
)

func main() {
	runtime.LockOSThread()
	//导入GWIN 主窗口控件
	gwin.ImportMainWindow() //可以在模板中使用MainWindow
	gwin.ImportChildWindow()
	gwin.ImportStaticText()
	//导入主模块描述
	bone.Import("AppMain", bone.Component{
		Template: `
		<MainWindow on:Destroy = "OnDestroy">
			<StaticText bind:Text="ChildText"></StaticText>
			<ChildWindow bind:_Bounds="ChildBounds,MultiWidth"></ChildWindow>
		</MainWindow>`,
		DataContext: reflect.TypeOf((*AppMainComponent)(nil)),
	}) //可以在模板中使用AppMain _开始的属性为单向数据绑定(没有_就是双向数据绑定)
	messageloop := gwin.NewMessageLoop()
	bootstrap := bone.NewBootstrap("<AppMain></AppMain>", nil, nil) //模板中使用AppMain
	messageloop.Append(bootstrap.Run)
	messageloop.EnterModal().LeaveModal().Run()
}

// AppMainComponent 主模块入口
type AppMainComponent struct {
	ChildBounds bone.DepProperty
	ChildText   bone.DepProperty
	exit        chan bool
}

func (am *AppMainComponent) Constructor(td bone.TaskDispatcher) {
	am.exit = make(chan bool, 1)
	am.ChildBounds = bone.NewDepProperty(200, nil, nil)
	time.AfterFunc(30*time.Second, func() {
		td.Invoke(func() {
			am.ChildBounds.Set(300)
		})
	})
	am.ChildText = bone.NewDepProperty(time.Now().Format("15:04:05"), nil, nil)
	go func() {
	loop:
		for {
			select {
			case <-am.exit:
				break loop
			case <-time.After(time.Second):
				td.Invoke(func() {
					am.ChildText.Set(time.Now().Format("15:04:05"))
				})
			}
		}
		fmt.Println("exit go routine success")
	}()

}

//MultiWidth ChildWindow 的bounds 是win.RECT 类型 而绑定上去的是整形 可以通过转换函数来返回win.RECT
//转换函数的设计 目的是为了支持 类似vue的计算属性
func (am *AppMainComponent) MultiWidth(sender interface{}, new interface{}) interface{} {
	rect := new.(int)
	return win.RECT{Left: 200, Top: 250, Right: int32(rect) * 2, Bottom: int32(rect) * 2}
}

//OnDestroy 主窗口被关闭的回调
func (am *AppMainComponent) OnDestroy(sender interface{}) {
	am.exit <- true
	win.PostQuitMessage(0)
}
