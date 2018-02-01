package bone

//Root 包含根节点
type Root struct {
	*View
}

func NewRoot(v *View) *Root {
	return &Root{
		View: v,
	}
}
