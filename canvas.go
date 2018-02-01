package bone

//Canvas 内置标签
type Canvas struct {
	*View
}

//NewCanvas 返回1个Canvas
func NewCanvas() *Canvas {
	return &Canvas{
		View: NewView("Canvas"),
	}
}
