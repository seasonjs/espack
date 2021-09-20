package config

// EntryPoints https://webpack.docschina.org/concepts/entry-points/
type EntryPoints interface{}

// EntryPointsResults 统一将类型转换为Map类型
type EntryPointsResults map[string]string

// OutputPoints  https://webpack.docschina.org/concepts/output/#root
type OutputPoints interface{}

// OutputPointsResults 统一将类型转换为结构体
type OutputPointsResults struct {
	Filename   string
	Path       string
	PublicPath string
}

// PluginPoints https://webpack.docschina.org/concepts/plugins/#anatomy
type PluginPoints interface{}

type configuration struct {
	// string
	Entry EntryPoints `json:"entry"`
	//string []string
	Output OutputPoints `json:"output"`

	Plugin PluginPoints `json:"plugin"`
}

//
func (c *configuration) SetEntry(entry EntryPoints) {
	c.Entry = entry
}

func (c *configuration) GetEntry(entry EntryPoints) {
	switch entry.(type) {
	case *string:
		c.Entry = entry.(*string)
	default:
		c.Entry = nil
	}

}

func (c *configuration) SetOutput(output OutputPoints) {
	c.Output = output
}

func (c *configuration) GetOutput(output OutputPoints) {
	switch output.(type) {
	case *string:
		c.Entry = output.(*string)
	default:
		c.Entry = nil
	}
}

func (c *configuration) SetPlugin(plugin PluginPoints) {
	c.Plugin = plugin
}

// GetPlugin TODO 转换插件类型
func (c *configuration) GetPlugin(plugin PluginPoints) {

}
