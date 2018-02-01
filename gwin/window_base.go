package gwin

import (
	"errors"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

var handle2wb = make(map[win.HWND]*windowBase)

type windowBase struct {
	Handle win.HWND
	Name   string

	proc func(msg uint32, wParam, lParam uintptr) (result uintptr)
}

func (wb *windowBase) Valid() bool {
	return wb.Handle != 0
}

func (wb *windowBase) Close() {
	if wb.Valid() {
		win.DestroyWindow(wb.Handle)
	}
}

func (wb *windowBase) wndProc(msg uint32, wParam, lParam uintptr) (result uintptr) {
	defer func() {
		if msg == win.WM_DESTROY {
			delete(handle2wb, wb.Handle)
			wb.Handle = 0
		}
	}()

	if wb.proc != nil {
		return wb.proc(msg, wParam, lParam)
	}

	return win.DefWindowProc(wb.Handle, msg, wParam, lParam)
}

func createWindowBase(parent *windowBase, className string, windowName string, style, exStyle uint32,
	x, y, nWidth, nHeight int32,
	wndProc func(msg uint32, wParam, lParam uintptr) (result uintptr)) (*windowBase, error) {
	parentHandle := win.HWND(0)
	if parent != nil {
		parentHandle = parent.Handle
	}
	wb := &windowBase{
		Name: windowName,
		proc: wndProc,
	}

	hwnd := win.CreateWindowEx(
		exStyle,
		syscall.StringToUTF16Ptr(className),
		syscall.StringToUTF16Ptr(windowName),
		style|win.WS_CLIPSIBLINGS,
		x,
		y,
		nWidth,
		nHeight,
		parentHandle,
		0,
		0,
		nil)
	if hwnd == 0 {
		return nil, errors.New("createWindow fail")
	}
	wb.Handle = hwnd
	handle2wb[hwnd] = wb
	return wb, nil
}

var (
	defaultWndProcPtr = syscall.NewCallback(defaultWndProc)
)

func defaultWndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	wb := handle2wb[hwnd]
	if wb != nil {
		return wb.wndProc(msg, wParam, lParam)
	}
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

func MustRegisterWindowClass(className string) {
	MustRegisterWindowClassWithWndProcPtr(className, defaultWndProcPtr)
}

func MustRegisterWindowClassWithStyle(className string, style uint32) {
	MustRegisterWindowClassWithWndProcPtrAndStyle(className, defaultWndProcPtr, style)
}

func MustRegisterWindowClassWithWndProcPtr(className string, wndProcPtr uintptr) {
	MustRegisterWindowClassWithWndProcPtrAndStyle(className, wndProcPtr, 0)
}

func MustRegisterWindowClassWithWndProcPtrAndStyle(className string, wndProcPtr uintptr, style uint32) {
	hInst := win.GetModuleHandle(nil)
	if hInst == 0 {
		panic("GetModuleHandle")
	}

	hIcon := win.LoadIcon(hInst, win.MAKEINTRESOURCE(7)) // rsrc uses 7 for app icon
	if hIcon == 0 {
		hIcon = win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_APPLICATION))
	}
	if hIcon == 0 {
		panic("LoadIcon")
	}

	hCursor := win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_ARROW))
	if hCursor == 0 {
		panic("LoadCursor")
	}

	var wc win.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.LpfnWndProc = wndProcPtr
	wc.HInstance = hInst
	wc.HIcon = hIcon
	wc.HCursor = hCursor
	wc.HbrBackground = win.COLOR_BTNFACE + 1
	wc.LpszClassName = syscall.StringToUTF16Ptr(className)
	wc.Style = style

	if atom := win.RegisterClassEx(&wc); atom == 0 {
		panic("RegisterClassEx")
	}
}
