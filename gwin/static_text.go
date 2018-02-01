package gwin

import (
	"reflect"
	"syscall"
	"unsafe"

	"github.com/dalixu/bone"
	"github.com/lxn/win"
)

//StaticTextComponent native child window 组件描述
var StaticTextComponent = bone.Component{
	Template: `
	<Canvas when:Attached="WhenViewAttached" when:Deattached = "WhenViewDeattached">
	</Canvas>
	`,
	DataContext: reflect.TypeOf((*StaticText)(nil)),
}

//ImportStaticText 导入StaticText标签
func ImportStaticText() error {
	return bone.Import("StaticText", StaticTextComponent)
}

//StaticText 静态文本
type StaticText struct {
	base   *windowBase
	logger bone.BLogger

	//依赖属性
	Widget  bone.DepProperty
	Bounds  bone.DepProperty //win.RECT
	Visible bone.DepProperty // bool
	Text    bone.DepProperty //string
}

func (st *StaticText) Constructor(logger bone.BLogger) {
	st.logger = logger
	//st.Bounds = bone.
	st.Bounds = bone.NewDepProperty(win.RECT{0, 0, 100, 200}, func(sender interface{}, value interface{}) {
		if st.base != nil && st.base.Valid() {
			bound := value.(win.RECT)
			if !win.SetWindowPos(st.base.Handle, 0, bound.Left, bound.Top, bound.Right-bound.Left, bound.Bottom-bound.Top, win.SWP_NOZORDER) {
				st.logger.Info("SetWindowPos fail")
			}
		}
	}, nil)
	st.Visible = bone.NewDepProperty(true, func(sender interface{}, value interface{}) {
		if st.base != nil && st.base.Valid() {
			visible := value.(bool)
			var show int32
			if visible {
				show = win.SW_SHOWNORMAL
			} else {
				show = win.SW_HIDE
			}
			if !win.ShowWindow(st.base.Handle, show) {
				st.logger.Info("ShowWindow fail")
			}
		}
	}, nil)
	st.Widget = bone.NewDepProperty((*bone.Widget)(nil), nil, nil)

	st.Text = bone.NewDepProperty("test", func(sender interface{}, value interface{}) {
		if st.base != nil && st.base.Valid() {
			text := value.(string)
			if win.TRUE != win.SendMessage(st.base.Handle, win.WM_SETTEXT, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text)))) {
				st.logger.Info("WM_SETTEXT failed")
			}
		}
	}, nil)
}

func (st *StaticText) wndProc(msg uint32, wParam, lParam uintptr) (result uintptr) {
	// if msg == win.WM_DESTROY {
	// 	if st.OnDestroy != nil {
	// 		st.OnDestroy(cw)
	// 	}
	// }
	return win.DefWindowProc(st.base.Handle, msg, wParam, lParam)
}

//WhenViewAttached 标签被添加到树上
func (st *StaticText) WhenViewAttached(sender *bone.View, parent, child interface{}) {
	if sender != child { //不是自己被添加到树上
		return
	}
	//子窗口必须要有父
	widget := sender.GetWidget()
	if widget == nil {
		st.logger.Error("StaticText must have a parent")
		return
	}
	//widget 设置到view
	var err error
	bound := st.Bounds.Get().(win.RECT)
	style := uint32(win.WS_CHILD | win.SS_CENTERIMAGE | win.WS_BORDER | win.SS_CENTER)
	if st.Visible.Get().(bool) {
		style |= win.WS_VISIBLE
	}
	st.base, err = createWindowBase(widget.Native.(*windowBase), "STATIC", st.Text.Get().(string),
		style, win.WS_EX_CONTROLPARENT, bound.Left, bound.Top,
		bound.Right-bound.Left, bound.Bottom-bound.Top, st.wndProc)
	if err != nil {
		st.logger.Error(err)
		return
	}
	st.Widget.Set(&bone.Widget{Native: st.base})
}

//WhenViewDeattached 标签从树上摘除
func (st *StaticText) WhenViewDeattached(sender *bone.View, parent, child interface{}) {
	if st.base != nil && st.base.Valid() {
		st.Widget.Set(nil)
		st.base.Close()
		st.base = nil
	}
}
