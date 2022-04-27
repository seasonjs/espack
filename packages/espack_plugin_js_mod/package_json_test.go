package espack_plugin_js_mod

import (
	"fmt"
	"seasonjs/espack/internal/utils"
	"testing"
)

func TestPackageJSON_ReadFile(t *testing.T) {
	path, _ := utils.FS.ConvertPath("../case/package.json")
	pkj := NewPackageJson().ReadFile(path).GetDependencies()
	fmt.Println(pkj)
}
