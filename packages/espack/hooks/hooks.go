package hooks

import (
	"github.com/evanw/esbuild/pkg/api"
	"github.com/seasonjs/espack/internal/config"
	"github.com/seasonjs/espack/internal/plugins"
	"time"
)

// HookContext Hook上下文
type hookContext struct {
	timer         time.Time
	pluginList    *plugins.PluginQueue  //插件队列，先进先出
	configuration *config.Configuration //配置
	result        *api.BuildResult
}
