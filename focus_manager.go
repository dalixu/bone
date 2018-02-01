package bone

type FocusTraversable interface {
	FocusSearch() FocusSearch
	FocusTraversableParent() FocusTraversable
	FocusTraversableParentView() View
}

type FocusChangeReason int

const (
	REASON_FOCUS_TRAVERSAL = iota
	REASON_FOCUS_STORE
	REASON_FOCUS_DIRECT_CHANGE
)

type FocusChangeDirection int

const (
	FOCUS_FORWARD = iota
	FOCUS_BACKWARD
)

type FocusManager struct {
}
