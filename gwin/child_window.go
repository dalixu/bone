package gwin

import (
	"reflect"

	"github.com/dalixu/bone"
	"github.com/lxn/win"
)

//ChildWindowComponent native child window 组件描述
var ChildWindowComponent = bone.Component{
	Template: `
	<Canvas when:Created = "WhenViewCreated" when:Attached="WhenViewAttached" when:Deattached = "WhenViewDeattached">
	</Canvas>
	`,
	DataContext: reflect.TypeOf((*ChildWindow)(nil)),
}

var childWindowClass = "gwin_ChildWindow"

//ImportChildWindow 导入子窗口标签
func ImportChildWindow() error {
	MustRegisterWindowClass(childWindowClass)
	return bone.Import("ChildWindow", ChildWindowComponent)
}

//ChildWindow native 窗口
type ChildWindow struct {
	base   *windowBase
	logger bone.BLogger
	widget *bone.Widget
	canvas *bone.View
	//依赖属性
	Bounds  bone.DepProperty //win.RECT
	Visible bone.DepProperty // bool
	//事件
	OnDestroy func(sender interface{})
}

//Constructor 构造函数
func (cw *ChildWindow) Constructor(logger bone.BLogger) {
	cw.logger = logger
	cw.Bounds = bone.NewDepProperty(win.RECT{0, 0, 100, 100}, func(sender interface{}, value interface{}) {
		if cw.base != nil && cw.base.Valid() {
			bound := value.(win.RECT)
			if !win.SetWindowPos(cw.base.Handle, 0, bound.Left, bound.Top, bound.Right-bound.Left, bound.Bottom-bound.Top, win.SWP_NOZORDER) {
				cw.logger.Info("SetWindowPos fail")
			}
		}
	}, nil)
	cw.Visible = bone.NewDepProperty(true, func(sender interface{}, value interface{}) {
		if cw.base != nil && cw.base.Valid() {
			visible := value.(bool)
			var show int32
			if visible {
				show = win.SW_SHOWNORMAL
			} else {
				show = win.SW_HIDE
			}
			if !win.ShowWindow(cw.base.Handle, show) {
				cw.logger.Info("ShowWindow fail")
			}
		}
	}, nil)
}

func (cw *ChildWindow) wndProc(msg uint32, wParam, lParam uintptr) (result uintptr) {
	if msg == win.WM_DESTROY {
		if cw.OnDestroy != nil {
			cw.OnDestroy(cw)
		}
	}
	return win.DefWindowProc(cw.base.Handle, msg, wParam, lParam)
}

//WhenViewCreated 内置标签view创建通知
func (cw *ChildWindow) WhenViewCreated(sender *bone.View) {
	cw.canvas = sender
}

//WhenViewAttached 标签被添加到树上
func (cw *ChildWindow) WhenViewAttached(sender *bone.View, parent, child interface{}) {
	if sender != child { //不是自己被添加到树上
		return
	}
	//子窗口必须要有父
	widget := sender.GetWidget()
	if widget == nil {
		cw.logger.Error("ChildWindow must have a parent")
		return
	}
	//widget 设置到view
	var err error
	bound := cw.Bounds.Get().(win.RECT)
	style := uint32(win.WS_CHILD)
	if cw.Visible.Get().(bool) {
		style |= win.WS_VISIBLE
	}
	cw.base, err = createWindowBase(widget.Native.(*windowBase), childWindowClass, "test",
		style, win.WS_EX_CONTROLPARENT, bound.Left, bound.Top,
		bound.Right-bound.Left, bound.Bottom-bound.Top, cw.wndProc)
	if err != nil {
		cw.logger.Error(err)
		return
	}
	cw.widget = &bone.Widget{Native: cw.base}
	cw.canvas.Widget.Set(cw.widget)
}

//WhenViewDeattached 标签从树上摘除
func (cw *ChildWindow) WhenViewDeattached(sender *bone.View, parent, child interface{}) {
	if cw.base != nil && cw.base.Valid() {
		cw.canvas.Widget.Set(nil)
		cw.base.Close()
		cw.base = nil
	}
}
