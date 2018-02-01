package bone

import (
	"fmt"
	"reflect"
	"strings"
)

func NewBone(ij *Injector, data interface{}, doc xmlDoc) *Bone {
	return &Bone{
		Injector: ij,
		data:     data,
		document: doc,
	}
}

type Bone struct {
	*Injector
	data interface{} //数据绑定上下文  通过反射自动绑定

	parent   *Bone
	children []*Bone
	topView  *View
	document xmlDoc
}

func (be *Bone) Parent() *Bone {
	return be.parent
}

func (be *Bone) ChildAt(index int) *Bone {
	if len(be.children) > index {
		return be.children[index]
	}
	return nil
}

func (be *Bone) RemoveChildren() {
	for len(be.children) > 0 {
		be.Remove(be.children[len(be.children)-1])
	}
}

func (be *Bone) Remove(child *Bone) {
	if child == nil || child.Parent() != be {
		return
	}
	index := -1
	for k, c := range be.children {
		if c == child {
			index = k
			break
		}
	}
	if index < 0 {
		return
	}
	child.parent = nil
	if index == len(be.children)-1 {
		be.children = be.children[0:index]
	} else if index == 0 {
		be.children = be.children[1:]
	} else {
		be.children = append(be.children[:index], be.children[index+1:]...)
	}
}

func (be *Bone) AddChild(child *Bone, index int) {
	if child == nil || be == child {
		return
	}
	//必须先被从树上remove掉才能添加
	if child.Parent() != nil || len(be.children) < index {
		return
	}
	if len(be.children) == index {
		be.children = append(be.children, child)
	} else {
		suffix := append([]*Bone(nil), be.children[index:]...)
		be.children = append(append(be.children[0:index], child), suffix...)
	}
	child.parent = be
}

func (be *Bone) AddChildToBack(child *Bone) {
	be.AddChild(child, len(be.children))
}

func (be *Bone) ChildrenCount() int {
	return len(be.children)
}

func (be *Bone) TopView() *View {
	if be.topView == nil {
		//要么 template 为空 要么他的template指向的是一个Bone 此时应该只有1个子
		if be.ChildrenCount() > 0 {
			return be.children[0].TopView()
		}
		return nil
	}
	return be.topView
}

func (be *Bone) load(parent *View) error {
	return be.create(be.document.Root, be, parent, CreateScope(be.data, "", nil))
}

// func (be *Bone) reload() error {
// 	var parentView *View
// 	index := -1
// 	view := be.TopView()
// 	if view != nil {
// 		parentView = view.Parent()
// 		if parentView != nil {
// 			index = parentView.IndexOfChild(view)
// 		}
// 	}
// 	//把bonechild 和topview从树上拿掉
// 	be.RemoveChildren()
// 	if parentView != nil {
// 		parentView.Remove(view)
// 	}

// 	if err := be.load(); err != nil {
// 		return err
// 	}
// 	// if be.hasIfForDirective(be.document.Root) {
// 	// 	return errors.New("top level can't have if or for directive")
// 	// }
// 	be.create(be.document.Root, be, nil, &DataContext{data: be.data})
// 	return nil
// }

func (be *Bone) create(node *xmlNode, thisBone *Bone, parentView *View, se Scope) error {
	//暂时不支持for if 等扩展指令
	//获得for指令
	// forDirective, err := be.ForDirective(node, dc)
	// if err != nil {
	// 	return err
	// }
	// ifDirecive, err := be.IfDirecive(node)
	// if err != nil {
	// 	return err
	// }
	count := 1
	// if forDirective != nil {
	// 	count = len(forDirective.DataContext)
	// }
	for ; count > 0; count-- {
		// create := true
		// if ifDirecive != nil {
		// 	create = false
		// }
		// if !create {
		// 	continue
		// }
		view := parentView
		bone := thisBone
		var err error
		//看node是否有for属性
		switch node.Element.Name.Local {
		case "Canvas":
			view, err = be.createCanvas(node, parentView, se)
			if err != nil {
				return err
			}

			if node == thisBone.document.Root {
				thisBone.topView = view
			}
		default:
			//是组件
			bone, err = be.createBone(node, thisBone, parentView)
			if err != nil {
				return err
			}
			//bone创建好后
			//执行数据绑定
			view = bone.TopView()
		}

		for _, v := range node.Child {
			if err = be.create(v, bone, view, se); err != nil {
				return err
			}
		}
	}
	return nil
}

func (be *Bone) createCanvas(node *xmlNode, parentView *View, se Scope) (*View, error) {
	var err error
	canvasView := NewCanvas().View
	for _, attr := range node.Element.Attr {
		name := attr.Name.Space
		if attr.Value == "" {
			continue
		}
		if name == "when" {
			err = bindNotifyOrEvent(canvasView, "When"+strings.TrimSpace(attr.Name.Local), strings.TrimSpace(attr.Value), se)
		} else if name == "on" { //事件绑定
			err = bindNotifyOrEvent(canvasView, "On"+strings.TrimSpace(attr.Name.Local), strings.TrimSpace(attr.Value), se)
		} else if name == "bind" { //数据绑定
			err = bindData(be.Get(reflect.TypeOf((*WatcherManager)(nil)).Elem()).Interface().(WatcherManager),
				canvasView, strings.TrimSpace(attr.Name.Local), strings.TrimSpace(attr.Value), se)
		}
		if err != nil {
			return nil, err
		}
	}
	//发送创建通知
	if canvasView.WhenCreated != nil {
		canvasView.WhenCreated(canvasView)
	}
	if parentView != nil {
		parentView.AddChildToBack(canvasView)
	}
	return canvasView, nil
}

func bindData(wm WatcherManager, view interface{}, name string, value string, se Scope) error {
	multi := true
	if name[0] == '_' { //带有"_"的是单向数据绑定 没有"_"的是双向数据绑定
		multi = false
		name = strings.TrimSpace(name[1:])
	}
	vName := name
	vValue := reflect.ValueOf(view)
	if vValue.Kind() == reflect.Ptr {
		vValue = vValue.Elem()
	}
	vProperty := vValue.FieldByName(vName)
	//必须是依赖属性
	if !vProperty.IsValid() {
		return fmt.Errorf("unknown method:%s", vName)
	}
	vmName := value
	vmFrom := ""
	vmTo := ""
	if values := strings.Split(value, ","); len(values) > 1 {
		vmName = strings.TrimSpace(values[0])
		vmFrom = strings.TrimSpace(values[1])
		if len(values) > 2 {
			vmTo = strings.TrimSpace(values[2])
		}
	}
	//先找func
	var vmComputerFrom func(sender interface{}, new interface{}) interface{}
	var vmComputerTo func(sender interface{}, new interface{}) interface{}
	if vmFrom != "" {
		fun, err := se.MethodByName(vmFrom, reflect.TypeOf((func(sender interface{}, new interface{}) interface{})(nil)))
		if err != nil {
			return err
		}
		vmComputerFrom = fun.(func(sender interface{}, new interface{}) interface{})
	}
	if vmTo != "" {
		fun, err := se.MethodByName(vmTo, reflect.TypeOf((func(sender interface{}, new interface{}) interface{})(nil)))
		if err != nil {
			return err
		}
		vmComputerTo = fun.(func(sender interface{}, new interface{}) interface{})
	}

	vmProperty, err := se.FieldByName(vmName, nil)
	if err != nil {
		return fmt.Errorf("unknown method:%s", vmName)
	}

	//找到了2个属性
	var sourceWatcher, targetWatcher *Watcher
	sourceWatcher = wm.Create(vmProperty.(WatchedTarget), vProperty.Interface().(WatchedTarget), vmComputerFrom)
	vmProperty.(WatchedTarget).AddSub(sourceWatcher)
	if multi {
		targetWatcher = wm.Create(vProperty.Interface().(WatchedTarget), vmProperty.(WatchedTarget), vmComputerTo)
		vProperty.Interface().(WatchedTarget).AddSub(targetWatcher)
	}
	sourceWatcher.Update() //同步一下属性
	return nil
}

func bindNotifyOrEvent(view interface{}, name string, value string, vm Scope) error {
	vName := name                   //target name
	vValue := reflect.ValueOf(view) //target value
	//view 以类似于C#的属性来声明 Method 变量
	if vValue.Kind() == reflect.Ptr {
		vValue = vValue.Elem()
	}
	vFunc := vValue.FieldByName(vName)
	if !vFunc.IsValid() {
		return fmt.Errorf("unknown method:%s", vName)
	}
	if !vFunc.CanSet() {
		return fmt.Errorf("view:CanSet false %+v", vFunc)
	}
	//
	vmFunc, err := vm.MethodByName(value, nil)
	if err != nil {
		return err
	}
	vFunc.Set(reflect.ValueOf(vmFunc))

	return nil
}

func (be *Bone) createBone(node *xmlNode, parent *Bone, parentView *View) (*Bone, error) {
	name := node.Element.Name.Local
	var ij *Injector
	if parent != nil {
		ij = parent.Injector
	}
	bone, err := createBoneByTag(name, ij)
	if err != nil {
		return nil, err
	}
	for _, attr := range node.Element.Attr {
		name := attr.Name.Space
		if attr.Value == "" {
			continue
		}
		if name == "when" {
			bindNotifyOrEvent(bone.data, "When"+strings.TrimSpace(attr.Name.Local), strings.TrimSpace(attr.Value), CreateScope(be.data, "", nil))
		} else if name == "on" { //事件绑定
			bindNotifyOrEvent(bone.data, "On"+strings.TrimSpace(attr.Name.Local), strings.TrimSpace(attr.Value), CreateScope(be.data, "", nil))
		} else if name == "bind" { //数据绑定
			bindData(bone.Get(reflect.TypeOf((*WatcherManager)(nil)).Elem()).Interface().(WatcherManager),
				bone.data, strings.TrimSpace(attr.Name.Local), strings.TrimSpace(attr.Value), CreateScope(be.data, "", nil))
		}
		if err != nil {
			return nil, err
		}
	}
	parent.AddChildToBack(bone) //bone先添加到父上
	err = bone.load(parentView)
	if err != nil {
		return nil, err
	}
	return bone, nil
}

// //节点是否有IF和FOR 指令
// func (be *Bone) hasIfForDirective(node *xmlNode) bool {
// 	for _, v := range node.Element.Attr {
// 		if v.Name.Local == "for" || v.Name.Local == "if" {
// 			return true
// 		}
// 	}
// 	return false
// }

// type IfDirecive struct {
// 	Key   string
// 	Value bool
// }

// type ForDirective struct {
// 	DataContext []interface{}
// 	Key         string //

// }

// func (be *Bone) IfDirecive(node *xmlNode) (*IfDirecive, error) {
// 	for _, v := range node.Element.Attr {
// 		if v.Name.Local == "if" {
// 			return &IfDirecive{
// 				Key: strings.TrimSpace(v.Value),
// 			}, nil
// 		}
// 	}
// 	return nil, nil
// }

// func (be *Bone) ForDirective(node *xmlNode, dc *DataContext) (*ForDirective, error) {
// 	for _, v := range node.Element.Attr {
// 		if v.Name.Local == "for" {
// 			//let name in dc
// 			//dc.FindField(v.Name.Local)
// 			return &ForDirective{}, nil
// 		}
// 	}
// 	return nil, nil
// }
