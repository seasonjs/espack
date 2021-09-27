package mod

import (
	"encoding/json"
	"fmt"
	"os"
	"seasonjs/espack/internal/utils"
)

// TODO 需要理解package.json每个字段的概念 https://docs.npmjs.com/cli/v7/configuring-npm/package-json
type packageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// NewPackageJson 每个包都要获取他的packages
func NewPackageJson() *packageJSON {
	return &packageJSON{}

}

// ReadFile 每个包都要获取他的packagejson
func (j *packageJSON) ReadFile(p ...string) *packageJSON {
	path := ""

	if p != nil && len(p[0]) > 0 {
		path = p[0]
	} else {
		path, _ = utils.FS.ConvertPath("./package.json")
	}

	file, _ := os.Open(path)
	// 关闭文件
	defer file.Close()
	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	decoder := json.NewDecoder(file)
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	err := decoder.Decode(j)
	if err != nil {
		fmt.Println("读取文件失败:", err)
	}
	return j
}

// GetDependencies 获取包依赖列表
func (j *packageJSON) GetDependencies() map[string]string {
	dependencies := make(map[string]string)
	for key, value := range j.Dependencies {
		dependencies[key] = value
	}
	//TODO: 需要移除包版本的前缀 例如: ^ ~
	for key, value := range j.DevDependencies {
		dependencies[key] = value
	}
	return dependencies
}
