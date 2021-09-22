package plugin

import "fmt"

func Register(name string, plugin func(opt interface{}) *interface{}) {
	//外部插件注册方案
	fmt.Printf("%s ,%v", name, plugin(nil))
}
