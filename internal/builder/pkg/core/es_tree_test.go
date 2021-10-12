package core

import (
	"testing"
)

func TestNewESTree(t *testing.T) {
	ast := NewESTree(VDefault, GDefault, JSX)

	t.Log(ast)
}
