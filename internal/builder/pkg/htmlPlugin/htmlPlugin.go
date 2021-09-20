package htmlPlugin

import (
	"github.com/evanw/esbuild/pkg/api"
	"io/ioutil"
	"seasonjs/espack/internal/plugins"
	"seasonjs/espack/internal/utils"
)

type htmlPlugin struct {
}

func NewHtmlPlugin() *htmlPlugin {
	return &htmlPlugin{}
}

func (p htmlPlugin) Setup() plugins.PluginResult {
	//TODO 替换为从配置中获得
	path, _ := utils.FS.ConvertPath("./index.html")
	outPath, _ := utils.FS.ConvertPath("./dist/index.html")
	buf, _ := ioutil.ReadFile(path)
	return plugins.PluginResult{
		PluginName: "espack_",
		OutputFile: api.OutputFile{
			Path:     outPath,
			Contents: buf,
		},
	}
}
