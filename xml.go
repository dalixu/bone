package bone

import (
	"bytes"
	"encoding/xml"
)

type xmlNode struct {
	Child   []*xmlNode
	Parent  *xmlNode
	Element xml.StartElement
}

//xmlDoc 把xml先解析成1棵树 然后遍历生成Bone 和View
type xmlDoc struct {
	Root *xmlNode
}

func (xd *xmlDoc) Parse(doc string) error {
	d := xml.NewDecoder(bytes.NewReader([]byte(doc)))
	token, err := d.Token()
	for {
		switch element := token.(type) {
		case xml.StartElement:
			root := &xmlNode{Element: element}
			err = recursiveLoad(root, d)
			if err == nil {
				xd.Root = root
			}
			return err
		default:
			token, err = d.Token()
			if err != nil {
				return err
			}
		}
	}
}

func recursiveLoad(node *xmlNode, d *xml.Decoder) error {
	for {
		token, err := d.Token()
		if err != nil {
			return err
		}
		switch element := token.(type) {
		case xml.StartElement:
			child := &xmlNode{Element: element}
			child.Parent = node
			node.Child = append(node.Child, child)
			if err = recursiveLoad(child, d); err != nil {
				return err
			}
		case xml.EndElement:
			return nil
		}
	}
}
