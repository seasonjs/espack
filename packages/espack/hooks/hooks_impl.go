package hooks

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/seasonjs/espack/internal/builder"
	"github.com/seasonjs/espack/internal/builder/pkg/html_plugin"
	"github.com/seasonjs/espack/internal/config"
	"github.com/seasonjs/espack/internal/devServer"
	"github.com/seasonjs/espack/internal/logger"
	"github.com/seasonjs/espack/internal/plugins"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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
	c.timer = time.Now()
	c.configuration = config.
		NewConfigPoints().
		ReadFile().
		ReadConfig()
	//创建plugin
	c.pluginList = plugins.NewPluginQueue()
	//TODO:通过文件引用插件，通过此处传入插件配置
	c.pluginList.Add(html_plugin.NewHtmlPlugin(html_plugin.HtmlPluginOption{}))
	return c
}

// InstallPlugin 执行插件,安装esbuild的 TODO:installPlugin
func (c *hookContext) InstallPlugin() *hookContext {
	//按照顺序调用
	for i := 0; i < c.pluginList.Len(); i++ {
		plugin := c.pluginList.Next()
		PluginResult := plugin.Setup(c.configuration)
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
		outputFiles := builder.NewEsBuilder(*c.configuration).GetOptions().EsbuildStarter().OutputFiles
		c.result.OutputFiles = append(c.result.OutputFiles, outputFiles...)
		expend := time.Since(c.timer)
		logger.Info("构建完成,共花费时间：%v", expend)
		buildFinish <- true
	}()
	return c
}

// LogAll 打印所有收集到的日志
func (c *hookContext) LogAll() *hookContext {
	return c
}

// HoldAll Hold主协程防止退出
func (c *hookContext) HoldAll() {
	sig := make(chan os.Signal)
	//监听所有信号
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//fmt.Println("start!")

	fmt.Println("接受到退出信号:", <-sig)
	os.Exit(1)
}
