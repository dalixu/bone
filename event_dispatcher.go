package bone

// //EventPath 事件传递路径
// type EventPath = *list.List

// //EventDispatcher 事件分发器
// type EventDispatcher struct {
// }

// //Push 分发1个Event
// func (ed EventDispatcher) Push(e Event) {
// 	if e == nil || e.Target() == nil {
// 		return
// 	}
// 	path := list.New()
// 	node := e.Target().Parent()
// 	for node != nil {
// 		path.PushBack(node)
// 		node = node.Parent()
// 	}
// 	e.setPhase(EP_CAPTURING)
// 	element := path.Front()
// 	for element != nil {
// 		element.Value.(View).handle(e)
// 		if !e.Propagation() {
// 			return
// 		}
// 		element = element.Next()
// 	}
// 	e.setPhase(EP_TARGET)
// 	e.Target().handle(e)
// 	if !e.Propagation() {
// 		return
// 	}
// 	if !e.Bubble() {
// 		return
// 	}
// 	e.setPhase(EP_BUBBLING)
// 	element = path.Back()
// 	for element != nil {
// 		element.Value.(View).handle(e)
// 		if !e.Propagation() {
// 			return
// 		}
// 		element = element.Prev()
// 	}
// }
