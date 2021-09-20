package hooks

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"os"
	"os/signal"
	"seasonjs/espack/internal/builder"
	"seasonjs/espack/internal/builder/pkg/htmlBuilder"
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
			result: &api.BuildResult{},
		}
		buildFinish = make(chan bool, 1)
	})

	return ctx
}

// InitHooks 初始化生命周期,做读取配置文件的操作并解析
func (c *hookContext) InitHooks() *hookContext {
	c.options = config.NewConfig().ReadFile().ReadConfig()
	//创建plugin
	c.pluginList = plugins.NewPluginQueue()
	//TODO:通过文件引用插件
	c.pluginList.Add(htmlBuilder.NewHtmlPlugin())
	return c
}

// InstallPlugin 执行插件,安装esbuild的 TODO:installPlugin
func (c *hookContext) InstallPlugin() *hookContext {
	//按照顺序调用
	for i := 0; i < c.pluginList.Len(); i++ {
		plugin := c.pluginList.Next()
		PluginResult := plugin.Setup()
		c.result.OutputFiles = append(c.result.OutputFiles, PluginResult.OutputFile)
	}
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
		//TODO 需要考虑被覆盖的问题
		c.result.OutputFiles = append(c.result.OutputFiles, builder.EsbuildStarter().OutputFiles...)
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
