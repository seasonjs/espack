package espack_plugin_js_mod

import (
	"fmt"
	"github.com/seasonjs/espack/internal/utils"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"
)

// 测试解压效果
func TestMod_UnTar(t *testing.T) {
	in, _ := utils.FS.ConvertPath("../case/yarn-0.1.0.tgz")
	out, _ := utils.FS.ConvertPath("../case/espack/mod")
	err := DefaultTarGz().UnTarGz(in, out)
	if err != nil {
		fmt.Println(err)
	}
}

const host = "https://registry.npmjs.org/"

// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

//req, err := http.NewRequest("GET", "https://registry.npmjs.org/@seasonjs/tools", nil)
//if err != nil {
//// handle err
//}
//req.Header.Set("Accept", "application/vnd.npm.install-v1+json")
//
//resp, err := http.DefaultClient.Do(req)
//if err != nil {
//// handle err
//}
//defer resp.Body.Close()
// 测试从npm获取全量元数据

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

func TestVersion_Reg(t *testing.T) {
	//v := versionReg.FindStringSubmatch("^1.1.0")
	//t.Log(v)
	var versionReg = regexp.MustCompile(`(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	match := versionReg.FindStringSubmatch("^1.1.0")

	results := map[string]string{}
	for i, name := range match {
		results[versionReg.SubexpNames()[i]] = name
	}
	t.Log(results)
	//return result
	//var myExp = regexp.MustCompile(`(?P<first>\d+)\.(\d+).(?P<second>\d+)`)
}

func TestDownloadFile(t *testing.T) {
	path, _ := utils.FS.ConvertPath("../case/node_modules/react")
	resp, _ := http.Get("https://registry.npmjs.org/react/-/react-17.0.2.tgz")
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		err := DefaultTarGz().IOUnTarGz(resp.Body, path)
		if err != nil {
			fmt.Println(err)
		}
		//bodyString := string(bodyBytes)
		fmt.Println("下载并解压成功")
	}
}
func TestDownloadFile2(t *testing.T) {
	path, _ := utils.FS.ConvertPath("../case/node_modules")
	jmP, _ := utils.FS.ConvertPath("../case/js.mod")
	filePath, _ := utils.FS.ConvertPath("../case/js.sum")
	NewMod().AnalyzeDependencies(jmP, filePath).DownLoadDependencies(path)
}