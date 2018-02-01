package bone

//Color 颜色
type Color uint32

//PMColor 预乘颜色
type PMColor uint32

//ColorF 颜色百分比
type ColorF struct {
	A float32
	R float32
	G float32
	B float32
}

//SetARGB 设置颜色
func (c *Color) SetARGB(a, r, g, b uint8) {
	*c = Color((uint32(a) << 24) | (uint32(r) << 16) | (uint32(g) << 8) | (uint32(b) << 0))
}

//GetA 得到alpha分量
func (c *Color) GetA() uint8 {
	return uint8((uint32(*c) >> 24) & 0xff)
}

//GetR 得到red分量
func (c *Color) GetR() uint8 {
	return uint8((uint32(*c) >> 16) & 0xff)
}

//GetG 得到green分量
func (c *Color) GetG() uint8 {
	return uint8((uint32(*c) >> 8) & 0xff)
}

//GetB 得到blue分量
func (c *Color) GetB() uint8 {
	return uint8((uint32(*c) >> 0) & 0xff)
}
