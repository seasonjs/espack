package config

//暂时不会支持webpack的复杂结构，支持的结构也需要进一步进行调研

// EntryPoints https://webpack.docschina.org/concepts/entry-points/
type EntryPoints interface{}

// EntryPointsResults 统一将类型转换为Map类型 TODO：支持复杂结构
type EntryPointsResults map[string]string

// ModeResults https://webpack.docschina.org/configuration/mode/
type ModeResults uint8

const (
	Production ModeResults = iota
	Development
	None
)

// OutputPoints https://webpack.docschina.org/configuration/output/
type OutputPoints interface{}

// OutputPointsResults 统一将类型转换为结构体 TODO：支持复杂结构
type OutputPointsResults struct {
	Filename   string
	Path       string
	PublicPath string
}

// PluginPoints https://webpack.docschina.org/concepts/plugins/#anatomy
type PluginPoints interface{}

// PluginPointsResults TODO： 调研插件类型结构
type PluginPointsResults struct {
}

// ConfigurationPoints JSON 类型声明
type ConfigurationPoints struct {
	//入口和上下文===============================================================
	//string https://webpack.docschina.org/configuration/entry-context/#context
	Context string
	// string | string[] | map[string]string
	Entry EntryPoints `json:"entry"`
	//模式=====================================================================
	//string = 'production': 'none' | 'development' | 'production'
	Mode string `json:"mode"`
	//string []string
	Output OutputPoints `json:"output"`

	Plugin PluginPoints `json:"plugin"`
}

// Configuration ConfigurationPoints 转换为 Configuration 固定类型声明
type Configuration struct {
	//string TODO:是否要考虑多包项目呢
	Context string
	//map[string]string
	Entry EntryPointsResults
	//string {}
	Mode ModeResults

	Output OutputPointsResults
	//TODO:确定配置方案
	Plugin PluginPointsResults
}

var cp ConfigurationPoints

func NewConfig(c ConfigurationPoints) *Configuration {
	cp = c
	return &Configuration{}
}

// 格式Context字段
func (c *Configuration) formatContext() *Configuration {
	c.Context = cp.Context
	return c
}

// 转换入口字段类型
func (c *Configuration) formatEntry() *Configuration {
	entry := make(EntryPointsResults)
	switch cp.Entry.(type) {
	case string:
		entry = EntryPointsResults{
			"main": cp.Entry.(string),
		}
	case []string:
		for i, s := range cp.Entry.([]string) {
			//TODO:需要转换成对应的名字
			entry[string(rune(i))] = s
		}
	case map[string]interface{}:
		for i, s := range cp.Entry.(map[string]interface{}) {
			entry[i] = s.(string)
		}
	default:
		//TODO:报错和提示
		entry = nil
	}
	c.Entry = entry
	return c
}

func (c *Configuration) formatMode() *Configuration {

	if cp.Mode == "development" {
		c.Mode = Development
		return c
	}
	if cp.Mode == "production" {
		c.Mode = Production
		return c
	}

	if cp.Mode == "none" {
		c.Mode = None
		return c
	}

	// 默认设置为 Production
	c.Mode = Production
	return c

}

// 转换输出类型
func (c *Configuration) formatOutput() *Configuration {
	output := OutputPointsResults{}
	switch cp.Output.(type) {
	case string:
		output = OutputPointsResults{
			Filename: "main",
			Path:     cp.Output.(string),
		}
	case OutputPointsResults:
		output = cp.Output.(OutputPointsResults)
	default:
		//TODO:报错和提示
		output = OutputPointsResults{}
	}
	c.Output = output
	return c
}

//TODO 转换插件类型
func (c *Configuration) formatPlugin() *Configuration {
	return c
}
