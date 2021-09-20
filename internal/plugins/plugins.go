package plugins

import (
	"github.com/evanw/esbuild/pkg/api"
	"sync"
)

type PluginResult struct {
	PluginName string
	OutputFile api.OutputFile
}

type Plugin interface {
	Setup() PluginResult //触发Plugin执行
}

// PluginQueue 插件队列
type PluginQueue struct {
	plugins []Plugin
	cur     int
	lock    sync.RWMutex
}
