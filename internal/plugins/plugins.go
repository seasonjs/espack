package plugins

import "sync"

type Plugin interface{}

// PluginQueue 插件队列
type PluginQueue struct {
	plugins []Plugin
	cur     int
	lock    sync.RWMutex
}
