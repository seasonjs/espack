package mod

import (
	"fmt"
	"io"
	"os"
	"seasonjs/espack/internal/utils"
	"strings"
	"testing"
)

func TestJsSum_FetchRootModMeta(t *testing.T) {
	path, _ := utils.FS.ConvertPath("../case/js.mod")
	require := NewJsMod().ReadFile(path).Require
	sum := NewJsSum().FetchRootModMeta(require)
	t.Log(sum)
}
func TestJsSum_PrunerAndFetcher(t *testing.T) {
	sum := NewJsSum()
	sum.modList = append(sum.modList, &modInfo{})
	sum.modList = nil
	sum.modList = append(sum.modList, &modInfo{})
	t.Log(sum.modList)
}
func TestJsSum_PrunerAndFetcher2(t *testing.T) {
	path, _ := utils.FS.ConvertPath("../case/js.mod")
	filePath, _ := utils.FS.ConvertPath("../case/js.sum")
	require := NewJsMod().ReadFile(path).Require
	js := NewJsSum().FetchRootModMeta(require).PrunerAndFetcher().WriteFile(filePath)
	t.Log(js)
}
func TestTrim(t *testing.T) {
	str := "@type/string@123456"
	for i := 0; i < 5; i++ {

	}
	str1 := strings.TrimRight(str, fmt.Sprintf("@%s", "123456"))
	t.Log(str1)
}
func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
func TestJsSum_WriteFile(t *testing.T) {
	//var writeString = "测试1\n测试2\n"
	filePath, _ := utils.FS.ConvertPath("../case/js.sum")
	// 覆盖写入，如果不存在就创建
	f, _ := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE, 0666) //打开文件
	defer f.Close()
	for i := 0; i < 5; i++ {
		_, _ = io.WriteString(f, fmt.Sprintf("xxxx %v ", i)) //写入文件(字符串)
	}

}
