package html_plugin

import (
	"github.com/evanw/esbuild/pkg/api"
	"github.com/seasonjs/espack/internal/config"
	"github.com/seasonjs/espack/internal/plugins"
	"github.com/seasonjs/espack/internal/utils"
	"io/ioutil"
)

// HtmlPluginOption https://github.com/jantimon/html-webpack-plugin#options TODO:考虑剩下的要支持哪些配置
type HtmlPluginOption struct {
	Title      string
	Filename   string
	Template   string
	PublicPath string
}

type htmlPlugin struct {
	opt HtmlPluginOption
}

func NewHtmlPlugin(opt HtmlPluginOption) *htmlPlugin {
	if len(opt.PublicPath) <= 0 {
		opt.PublicPath = "auto"
	}
	return &htmlPlugin{
		opt,
	}
}

func (p htmlPlugin) Setup(config *config.Configuration) plugins.PluginResult {
	//TODO 替换为从配置中获得
	path, _ := utils.FS.ConvertPath("./index.html")
	outPath, _ := utils.FS.ConvertPath("./dist/index.html")
	buf, _ := ioutil.ReadFile(path)
	return plugins.PluginResult{
		PluginName: "espack_html_plugin",
		OutputFile: api.OutputFile{
			Path:     outPath,
			Contents: buf,
		},
	}
}
