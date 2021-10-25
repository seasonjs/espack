package core

import (
	"testing"
)

type T interface {
	t()
}

type base struct {
}

func (b base) t() {

}

type A struct {
	base
}

func (b A) t() {

}

func NewBase() T {
	return base{}
}
func TestNewESTree(t *testing.T) {
	//ast := NewESTree(VDefault, GDefault, JSX)
	//
	//t.Log(ast)
	b := NewBase()
	a := A{
		b.(base),
	}
	t.Log(a)
}
