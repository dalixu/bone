package bone

import (
	"log"
	"testing"
)

func TestGet(t *testing.T) {
	vw := NewView()
	if vw.FirstChild() != nil {
		log.Fatal("")
	}
}
