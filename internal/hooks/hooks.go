package hooks

import (
	"github.com/evanw/esbuild/pkg/api"
	"seasonjs/espack/internal/config"
	"seasonjs/espack/internal/plugins"
)

// HookContext Hook上下文
type hookContext struct {
	pluginList    *plugins.PluginQueue  //插件队列，先进先出
	configuration *config.Configuration //配置
	result        *api.BuildResult
}
