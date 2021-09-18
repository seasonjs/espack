package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func NewConfig(arg ...interface{}) *configuration {
	if len(arg) > 1 {
		fmt.Println("New Config Error: Expect null or 1 arg , will return an empty configuration")
		return &configuration{}
	}
	if len(arg) < 1 {
		return &configuration{}
	}
	// len(arg) = 1
	if arg[0] != nil {
		switch arg[0].(type) {
		case *configuration:
			return arg[0].(*configuration)
		default:
			fmt.Println("New Config Error: arg type error, will return an empty configuration")
			return &configuration{}
		}
	}
	return &configuration{}
}

func (c *configuration) ReadFile(arg ...interface{}) *configuration {
	if len(arg) != 1 {
		fmt.Println("Read File Error: Expect 1 arg , will return an configuration without reading")
		return c
	}
	path := arg[0]
	if path == nil {
		// TODO: 替换为正确的路径
		path = "./espack.config.json"
	}
	file, _ := os.Open(fmt.Sprintf("%v", path))

	// 关闭文件
	defer file.Close()

	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	decoder := json.NewDecoder(file)
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	err := decoder.Decode(c)
	if err != nil {
		msg := fmt.Errorf("Read Config File Error: %v \n", err)
		panic(msg)
	}
	return c
}

//TODO
func (c *configuration) WatchFileChange() {

}

func (c *configuration) ReadConfig() *configuration {
	return c
}

//TODO: getter setter
