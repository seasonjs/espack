package plugins

// NewPluginQueue 构造函数创建插件队列
func NewPluginQueue() *PluginQueue {
	s := &PluginQueue{}
	s.cur = -1
	s.plugins = []Plugin{}
	return s
}

// Add 添加插件
func (p *PluginQueue) Add(t Plugin) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.plugins = append(p.plugins, t)
}

// Remove 移除插件
func (p *PluginQueue) Remove() Plugin {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.Len() == 0 {
		return nil
	}
	t := p.plugins[0]
	p.plugins = p.plugins[1:]
	return t
}

// Next 下一个插件
func (p *PluginQueue) Next() Plugin {
	if p.Len() == 0 {
		return nil
	}
	if p.cur >= p.Len()-1 {
		return nil
	}
	p.cur = p.cur + 1
	t := p.plugins[p.cur]
	return t
}

// Prev 上一个插件
func (p *PluginQueue) Prev() Plugin {
	if p.Len() == 0 {
		return nil
	}
	if p.cur <= 0 {
		return nil
	}
	p.cur = p.cur - 1
	t := p.plugins[p.cur]
	return t
}

// Len 插件长度
func (p *PluginQueue) Len() int {
	return len(p.plugins)
}
