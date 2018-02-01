package bone

//NewView 返回1个View go不能继承 所以其他内置标签只能包含1个view 为了区别每个view
//使用className
func NewView(className string) *View {
	v := &View{
		className: className,
		Widget:    NewDepProperty(nil, nil, nil),
	}
	return v
}

type View struct {
	className string

	parent   *View
	children []*View
	//依赖属性
	Widget DepProperty
	//通知
	WhenCreated    func(sender *View)
	WhenAttached   func(sender *View, parent, child interface{})
	WhenDeattached func(sender *View, parent, child interface{})
	//事件
}

//ClassName 返回类型
func (v *View) ClassName() string {
	return v.className
}

//GetWidget 获得Widget
func (v *View) GetWidget() *Widget {
	current := v
	for current != nil {
		if current.Widget.Get() != nil {
			return current.Widget.Get().(*Widget)
		}
		current = current.Parent()
	}
	return nil
}

func (v *View) Parent() *View {
	return v.parent
}

func (v *View) ChildAt(index int) *View {
	if len(v.children) > index {
		return v.children[index]
	}
	return nil
}

func (v *View) RemoveChildren() {
	for len(v.children) > 0 {
		v.Remove(v.children[len(v.children)-1])
	}
}

func (v *View) Remove(child *View) {
	if child == nil || child.Parent() != v {
		return
	}
	index := -1
	for k, c := range v.children {
		if c == child {
			index = k
			break
		}
	}
	if index < 0 {
		return
	}
	child.parent = nil
	if index == len(v.children)-1 {
		v.children = v.children[0:index]
	} else if index == 0 {
		v.children = v.children[1:]
	} else {
		v.children = append(v.children[:index], v.children[index+1:]...)
	}
	child.recursionWhenDeattached(v, child)
}

func (v *View) AddChild(child *View, index int) {
	if child == nil || v == child {
		return
	}
	//必须先被从树上remove掉才能添加
	if child.Parent() != nil || len(v.children) < index {
		return
	}
	if len(v.children) == index {
		v.children = append(v.children, child)
	} else {
		suffix := append([]*View(nil), v.children[index:]...)
		v.children = append(append(v.children[0:index], child), suffix...)
	}
	child.parent = v

	child.recursionWhenAttached(v, child)
}

func (v *View) recursionWhenAttached(parent, child *View) {
	if child == nil {
		return
	}
	//往下通知
	if child.WhenAttached != nil {
		child.WhenAttached(child, parent, child)
	}

	for _, v := range child.children {
		v.recursionWhenAttached(parent, child)
	}
}

func (v *View) recursionWhenDeattached(parent, child *View) {
	if child == nil {
		return
	}
	//往下通知
	if child.WhenDeattached != nil {
		child.WhenDeattached(child, parent, child)
	}

	for _, v := range child.children {
		v.recursionWhenDeattached(parent, child)
	}
}

func (v *View) IndexOfChild(child *View) int {
	if child == nil || child.Parent() != v {
		return -1
	}
	for k, c := range v.children {
		if c == child {
			return k
		}
	}
	return -1
}

func (v *View) AddChildToBack(child *View) {
	v.AddChild(child, len(v.children))
}

func (v *View) ChildrenCount() int {
	return len(v.children)
}

func (v *View) handle(e Event) {

}
