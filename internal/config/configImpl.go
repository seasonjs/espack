package config

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"seasonjs/espack/internal/utils"
)

func NewConfig(arg ...interface{}) *configuration {
	var length = len(arg)
	if length == 1 && arg[0] != nil {
		switch arg[0].(type) {
		case *configuration:
			return arg[0].(*configuration)
		default:
			return &configuration{}
		}
	} else {
		return &configuration{}
	}
}

func (c *configuration) ReadFile(arg ...interface{}) *configuration {
	var path = ""
	var msg = "Read Config Args Error: %s , will use default configuration "
	var length = len(arg)
	if length == 1 && arg[0] != nil {
		switch arg[0].(type) {
		case *string:
			path = arg[0].(string)
		default:
			fmt.Printf(msg, "args type unexpected")
			path, _ = utils.FS.ConvertPath("./espack.config.json")
		}
	} else {
		fmt.Printf(msg, "args length unexpected")
		path, _ = utils.FS.ConvertPath("./espack.config.json")
		fmt.Println(path)
	}

	file, _ := os.Open(path)
	// 关闭文件
	defer file.Close()

	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	decoder := json.NewDecoder(file)
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	err := decoder.Decode(c)
	if err != nil {
		utils.Err.LogAndExit(errors.Wrap(err, "Read Config File Error"))
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
