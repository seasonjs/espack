package mod

import (
	"fmt"
	"seasonjs/espack/internal/utils"
	"testing"
)

func TestJsMod_ReadFile(t *testing.T) {
	path, _ := utils.FS.ConvertPath("../case/js.mod")
	pkj := NewJsMod().ReadFile(path)
	fmt.Println(pkj)
}
