package mod

import (
	"fmt"
	"github.com/mholt/archiver/v3"
	"io/ioutil"
	"net/http"
	"seasonjs/espack/internal/utils"
	"testing"
)

// 测试解压效果
func TestMod_UnTar(t *testing.T) {

	in, _ := utils.FS.ConvertPath("../case/yarn-1.22.11.tgz")
	out, _ := utils.FS.ConvertPath("../case/espack/mod")
	err := archiver.Unarchive(in, out)
	if err != nil {
		fmt.Println(err)
	}
}

const host = "https://registry.npmjs.org/"

// 测试从npm获取元数据
func TestMod_Fetch(t *testing.T) {
	req, _ := http.NewRequest("GET", host+"@seasonjs/tools", nil)
	// 通过设置请求头缩小元数据量
	req.Header.Set("Accept", "application/vnd.npm.install-v1+json")
	resp, err := (&http.Client{}).Do(req)
	//resp, err := http.Get(serviceUrl + "/topic/query/false/lsj")
	if err != nil {
		fmt.Println("请求失败", err.Error())
		t.FailNow()
	}
	defer resp.Body.Close()
	s, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(s))
}
