package gwin

import (
	"time"

	"github.com/lxn/win"
)

//MessageLoop 消息循环 目前不能兼容模态对话框
type MessageLoop struct {
	idle    []func()
	exit    chan int
	modal   bool
	message *windowBase
	timer   uintptr
}

func init() {
	MustRegisterWindowClass("gwin_message")
}

//NewMessageLoop 返回消息循环
func NewMessageLoop() *MessageLoop {
	ml := &MessageLoop{
		exit:  make(chan int, 1),
		modal: false,
	}
	window, err := createWindowBase(&windowBase{Handle: win.HWND_MESSAGE}, "gwin_message", "", 0, 0, 0, 0, 0, 0, ml.wndProc)
	if err != nil {
		return nil
	}
	ml.message = window
	return ml
}

func (ml *MessageLoop) Close() {
	ml.message.Close()
}

func (ml *MessageLoop) wndProc(msg uint32, wParam, lParam uintptr) (result uintptr) {
	if msg == win.WM_TIMER {
		ml.idleLoop()
	}
	return win.DefWindowProc(ml.message.Handle, msg, wParam, lParam)
}

//Append 添加一个空闲例程
func (ml *MessageLoop) Append(idle func()) *MessageLoop {
	ml.idle = append(ml.idle, idle)
	return ml
}

//Quit 退出当前消息循环 必须LeaveModal后才能调用
func (ml *MessageLoop) Quit(code int) *MessageLoop {
	ml.exit <- code
	return ml
}

//Run 消息主循环
func (ml *MessageLoop) Run() int {
	var msg win.MSG
	for {
		select {
		case <-time.After(10 * time.Millisecond):
			ml.idleLoop()
			if win.PeekMessage(&msg, 0, 0, 0, win.PM_REMOVE) {
				if msg.Message == win.WM_QUIT {
					return int(msg.WParam)
				}
				win.TranslateMessage(&msg)
				win.DispatchMessage(&msg)
			}
		case ret := <-ml.exit:
			return ret
		}

		//runtime.Gosched()cpu 占用过高
	}
}

func (ml *MessageLoop) idleLoop() {
	for _, v := range ml.idle {
		v()
	}
}

//EnterModal 显示native dialog前 调用EnterModal()
// EnterModal->MessageBox->MessageBox->LeaveModal
func (ml *MessageLoop) EnterModal() *MessageLoop {
	if ml.modal {
		return ml
	}
	ml.modal = true
	win.SetTimer(ml.message.Handle, 0, 1000, 0)
	return ml
}

//LeaveModal 关闭native dialog后 调用
func (ml *MessageLoop) LeaveModal() *MessageLoop {
	if !ml.modal {
		return ml
	}
	win.KillTimer(ml.message.Handle, 0)
	ml.modal = false
	return ml
}
