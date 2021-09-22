package plugins

import (
	"github.com/evanw/esbuild/pkg/api"
	"seasonjs/espack/internal/config"
	"sync"
)

type PluginResult struct {
	PluginName string
	OutputFile api.OutputFile
}

type Plugin interface {
	// Setup 不可以让插件更改到配置信息
	Setup(points *config.Configuration) PluginResult //触发Plugin执行
}

// PluginQueue 插件队列
type PluginQueue struct {
	plugins []Plugin
	cur     int
	lock    sync.RWMutex
}
