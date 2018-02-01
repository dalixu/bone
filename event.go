package bone

//EventType 事件类型
type EventType int

//鼠标事件类型
const (
	ET_MOUSE_ENTER = iota
	ET_MOUSE_MOVE
	ET_MOUSE_DOWN
	ET_MOUSE_UP
	ET_MOUSE_CLICK
	ET_MOUSE_DCLICK
	ET_MOUSE_LEAVE
	ET_MOUSE_WHEEL

	ET_KEY_DOWN
	ET_CHAR
	ET_KEY_UP

	ET_COMPOSITION_START
	ET_COMPOSITION_UPDATE
	ET_COMPOSITION_END

	ET_FOCUS_OUT
	ET_FOCUS_IN
	ET_BLUR
	ET_FOCUS

	ET_CUSTOM
	ET_COUNT
)

//EventFlags like cef
type EventFlags int

//EventFlags定义
const (
	EF_NONE              = 0 // Used to denote no flags explicitly
	EF_CAPS_LOCK_ON      = 1 << 0
	EF_SHIFT_DOWN        = 1 << 1
	EF_CONTROL_DOWN      = 1 << 2
	EF_ALT_DOWN          = 1 << 3
	EF_LEFT_MOUSE_DOWN   = 1 << 4
	EF_MIDDLE_MOUSE_DOWN = 1 << 5
	EF_RIGHT_MOUSE_DOWN  = 1 << 6
	EF_COMMAND_DOWN      = 1 << 7 // GUI Key (e.g. Command on OS X keyboards,
	// Search on Chromebook keyboards,
	// Windows on MS-oriented keyboards)
	EF_NUM_LOCK_ON = 1 << 8
	EF_IS_KEY_PAD  = 1 << 9
	EF_IS_LEFT     = 1 << 10
	EF_IS_RIGHT    = 1 << 11
)

//EventPhase 事件阶段
type EventPhase int

//EventPhase 常量定义
const (
	EP_CAPTURING = 1 + iota
	EP_TARGET
	EP_BUBBLING
)

//Event 鼠标键盘及自定义事件对象
type Event interface {
	Target() View
	Propagation() bool
	Bubble() bool

	setPhase(phase EventPhase)
}

type event struct {
	et          EventType
	target      View
	phase       EventPhase
	bubble      bool
	cancelable  bool
	cancelled   bool
	propagation bool
}

func newEvent() Event {
	return &event{}
}

func (e *event) StopPropagation() {
	e.propagation = false
}

func (e *event) PreventDefault() {
	if e.cancelable {
		e.cancelled = true
	}
}

func (e *event) Cancelled() bool {
	if e.cancelable {
		return e.cancelled
	}
	return false
}

func (e *event) Target() View {
	return e.target
}

func (e *event) Propagation() bool {
	return e.propagation
}
func (e *event) Bubble() bool {
	return e.bubble
}

func (e *event) Type() EventType {
	return e.et
}

func (e *event) Phase() EventPhase {
	return e.phase
}

func (e *event) setPhase(phase EventPhase) {
	e.phase = phase
}
