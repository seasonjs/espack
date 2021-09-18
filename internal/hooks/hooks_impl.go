package hooks

import (
	"fmt"
	"os"
	"os/signal"
	"seasonjs/espack/internal/builder"
	"seasonjs/espack/internal/config"
	"seasonjs/espack/internal/devServer"
	"seasonjs/espack/internal/plugins"
	"sync"
)

var (
	ctx         *hookContext
	once        sync.Once
	buildFinish chan bool // esbuild 是否完成构建
)

func NewHookContext() *hookContext {
	once.Do(func() {
		ctx = &hookContext{
			pluginList: plugins.NewPluginQueue(),
			options:    nil,
		}
		buildFinish = make(chan bool, 1)
	})

	return ctx
}

// InitHooks 初始化生命周期,做读取配置文件的操作
func (c *hookContext) InitHooks() *hookContext {
	c.options = config.NewConfig().ReadFile().ReadConfig()
	return c
}

// InstallPlugin 安装插件,目前是安装esbuild的
func (c *hookContext) InstallPlugin() *hookContext {
	//创建plugin
	c.pluginList = plugins.NewPluginQueue()
	return c
}

// StartDevServer  开始启动dev如果是生成环境则不需要启动
func (c *hookContext) StartDevServer() *hookContext {
	ctx := devServer.NewDevServer()
	ctx.Run()
	go func() {
		//esbuild 完成才能继续执行
		if <-buildFinish {
			ctx.Add(&c.result.OutputFiles)
		}
	}()
	return c
}

// StartESBuild 启动Esbuild
func (c *hookContext) StartESBuild() *hookContext {

	go func() {
		c.result = builder.EsbuildStarter()
		buildFinish <- true
		fmt.Println("esbuild finish")
	}()
	return c
}
func (c *hookContext) HoldAll() {
	sig := make(chan os.Signal)
	//监听所有信号
	signal.Notify(sig)
	fmt.Println("start!")
	fmt.Println("stop,signal : ", <-sig)
}
