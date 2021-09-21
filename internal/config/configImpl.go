package config

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"seasonjs/espack/internal/utils"
)

func NewConfigPoints(arg ...interface{}) *ConfigurationPoints {
	var length = len(arg)
	if length == 1 && arg[0] != nil {
		switch arg[0].(type) {
		case *ConfigurationPoints:
			return arg[0].(*ConfigurationPoints)
		default:
			return &ConfigurationPoints{}
		}
	} else {
		return &ConfigurationPoints{}
	}
}

// ReadFile TODO 读取配置文件
func (c *ConfigurationPoints) ReadFile(arg ...interface{}) *ConfigurationPoints {
	var path = ""
	var length = len(arg)
	if length == 1 && arg[0] != nil {
		switch arg[0].(type) {
		case string:
			path = arg[0].(string)
		default:
			path, _ = utils.FS.ConvertPath("./espack.config.json")
		}
	} else {
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

// WatchFileChange TODO
func (c *ConfigurationPoints) WatchFileChange() {

}

// ReadConfig 在此处进行配置的读取和转换 注意！！！ 必须在执行完成 ReadFile后调用
func (c *ConfigurationPoints) ReadConfig() *Configuration {
	config := NewConfig(*c).formatEntry().formatOutput().formatPlugin()
	return config
}
