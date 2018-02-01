package gwin

import (
	"reflect"

	"github.com/lxn/win"

	"github.com/dalixu/bone"
)

//MainWindowComponent Native MainWindow 组件描述
var MainWindowComponent = bone.Component{
	Template: `
	<Canvas bind:Widget="Widget" when:Created = "WhenViewCreated" when:Attached="WhenViewAttached" when:Deattached = "WhenViewDeattached">
	</Canvas>
	`,
	DataContext: reflect.TypeOf((*MainWindow)(nil)),
}

var mainWindowClass = "gwin_MainWindow"

//ImportMainWindow 导入主窗口标签
func ImportMainWindow() error {
	MustRegisterWindowClass(mainWindowClass)
	return bone.Import("MainWindow", MainWindowComponent)
}

//MainWindow native 窗口
type MainWindow struct {
	base   *windowBase
	logger bone.BLogger
	canvas *bone.View

	//依赖属性 因为要通过反射来绑定 必须首字母大写
	//Widget也许不需要是依赖属性 因为大写导出了 可能会被其他标签绑定
	//要么使用canvas直接设置值 要么后期对DepProperty增加私有权限控制
	Widget bone.DepProperty //*bone.Widget

	Bounds  bone.DepProperty //win.RECT
	Visible bone.DepProperty // bool
	//事件
	OnDestroy func(sender interface{})
}

//Constructor 构造函数 顶级窗口可以直接创建
func (mw *MainWindow) Constructor(logger bone.BLogger) {
	mw.logger = logger
	mw.Bounds = bone.NewDepProperty(win.RECT{0, 0, 1024, 1024}, func(sender interface{}, value interface{}) {
		if mw.base != nil && mw.base.Valid() {
			bound := value.(win.RECT)
			if !win.SetWindowPos(mw.base.Handle, 0, bound.Left, bound.Top, bound.Right-bound.Left, bound.Bottom-bound.Top, win.SWP_NOZORDER) {
				mw.logger.Info("SetWindowPos fail")
			}
		}
	}, nil)
	mw.Visible = bone.NewDepProperty(true, func(sender interface{}, value interface{}) {
		if mw.base != nil && mw.base.Valid() {
			visible := value.(bool)
			var show int32
			if visible {
				show = win.SW_SHOWNORMAL
			} else {
				show = win.SW_HIDE
			}
			if !win.ShowWindow(mw.base.Handle, show) {
				mw.logger.Info("ShowWindow fail")
			}
		}
	}, nil)
	mw.Widget = bone.NewDepProperty((*bone.Widget)(nil), nil, nil)
}

func (mw *MainWindow) wndProc(msg uint32, wParam, lParam uintptr) (result uintptr) {
	if msg == win.WM_DESTROY {
		if mw.OnDestroy != nil {
			mw.OnDestroy(mw)
		}
	}
	return win.DefWindowProc(mw.base.Handle, msg, wParam, lParam)
}

//WhenViewCreated 内置标签view创建通知
func (mw *MainWindow) WhenViewCreated(sender *bone.View) {
	mw.canvas = sender
}

//WhenViewAttached 标签被添加到树上
func (mw *MainWindow) WhenViewAttached(sender *bone.View, parent, child interface{}) {
	if sender != child { //不是自己被添加到树上
		return
	}
	//widget 设置到view上
	var err error
	style := uint32(win.WS_OVERLAPPEDWINDOW)
	if mw.Visible.Get().(bool) {
		style |= win.WS_VISIBLE
	}
	bound := mw.Bounds.Get().(win.RECT)
	mw.base, err = createWindowBase(nil, mainWindowClass, "", style, win.WS_EX_CONTROLPARENT,
		bound.Left, bound.Top,
		bound.Right-bound.Left, bound.Bottom-bound.Top, mw.wndProc)
	if err != nil {
		mw.logger.Error(err)
		return
	}
	mw.Widget.Set(&bone.Widget{Native: mw.base})
	//mw.canvas.Widget.Set(nil, mw.widget)
}

//WhenViewDeattached 标签从树上摘除
func (mw *MainWindow) WhenViewDeattached(sender *bone.View, parent, child interface{}) {
	if mw.base != nil && mw.base.Valid() {
		mw.Widget.Set((*bone.Widget)(nil))
		//mw.canvas.Widget.Set(nil, nil)
		mw.base.Close()
		mw.base = nil
	}
}
