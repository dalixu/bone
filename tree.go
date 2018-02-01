package bone

type treeNode struct {
	value         interface{}
	parent        *treeNode
	firstChild    *treeNode
	prevSibling   *treeNode
	nextSibling   *treeNode
	childrenCount int
}

func newTreeNode(value interface{}) *treeNode {
	return &treeNode{
		value:         value,
		parent:        nil,
		firstChild:    nil,
		prevSibling:   nil,
		nextSibling:   nil,
		childrenCount: 0,
	}
}
func (tn *treeNode) Value() interface{} {
	return tn.value
}

func (tn *treeNode) Parent() *treeNode {
	return tn.parent
}

func (tn *treeNode) FirstChild() *treeNode {
	return tn.firstChild
}

func (tn *treeNode) LastChild() *treeNode {
	if tn.firstChild != nil {
		return tn.firstChild.prevSibling
	}
	return nil
}

func (tn *treeNode) PrevSibling() *treeNode {
	if tn.parent != nil && tn.parent.firstChild != tn {
		return tn.prevSibling
	}
	return nil
}

func (tn *treeNode) NextSibling() *treeNode {
	if tn.parent != nil && tn.parent.firstChild != tn.nextSibling {
		return tn.nextSibling
	}
	return nil
}

func (tn *treeNode) Remove(child *treeNode) {
	if child == nil || child == tn {
		return
	}
	if child.parent != tn {
		return
	}
	if 1 == tn.childrenCount { //只有1个子
		tn.firstChild = nil
	} else { //多个子
		if child == tn.firstChild {
			tn.firstChild = child.nextSibling
		}
		child.prevSibling.nextSibling = child.nextSibling
		child.nextSibling.prevSibling = child.prevSibling
	}
	tn.childrenCount--
	child.parent = nil
}

func (tn *treeNode) RemoveFromParent() {
	if tn.parent != nil {
		tn.parent.Remove(tn)
	}
}

func (tn *treeNode) AddChild(child *treeNode, index int) {
	if child == nil || tn == child {
		return
	}
	//必须先被从树上remove掉才能添加
	if child.parent != nil || tn.childrenCount < index {
		return
	}
	if tn.firstChild != nil { //有子
		old := tn.firstChild
		if index < tn.childrenCount { //在中间插入
			for i := index; i > 0; i-- {
				old = old.nextSibling
			}
		}
		child.nextSibling = old
		child.prevSibling = old.prevSibling
		old.prevSibling.nextSibling = child
		old.prevSibling = child
		if 0 == index {
			tn.firstChild = child
		}
	} else { //没有子
		child.nextSibling = child
		child.prevSibling = child
		tn.firstChild = child
	}
	tn.childrenCount++
	child.parent = tn
}

func (tn *treeNode) AddChildToBack(child *treeNode) {
	tn.AddChild(child, tn.childrenCount)
}

func (tn *treeNode) ChildrenCount() int {
	return tn.childrenCount
}
