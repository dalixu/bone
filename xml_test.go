package bone

import "testing"
import "log"

type A interface {
	Add()
}

type B interface {
	Sub()
}

type MA struct {
}

func (ma *MA) Add() {

}

func (ma *MA) Sub() {

}

func TestParse(t *testing.T) {
	ma := new(MA)
	a := interface{}(ma)
	b := A(ma)
	c := B(ma)
	log.Printf("%p\n", ma)
	log.Printf("%p\n", a)
	log.Printf("%p\n", b)
	log.Printf("%p\n", c)
	doc := (&xmlDoc{})
	doc.Parse("<a><b><d></d></b><c></c></a>")
	for _, v := range doc.Root.Child {
		log.Printf("%+v\n", v.Element.Name)
	}
}
