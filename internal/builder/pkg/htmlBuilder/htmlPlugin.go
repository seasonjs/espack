package htmlBuilder

import (
	"github.com/evanw/esbuild/pkg/api"
	"io/ioutil"
	"seasonjs/espack/internal/utils"
)

//func NewHtmlPlugin() api.Plugin {
//	return api.Plugin{
//		Name: "espack_html_plugin",
//		Setup: func(build api.PluginBuild) {
//			//build.OnStart= func(callback func() (api.OnStartResult, error)) {
//			//
//			//}
//			build.OnLoad = loadHtml
//			build.OnResolve = resolveHtmlPath
//		},
//	}
//}

// 找到js 引入的html文件路径
func resolveHtmlPath(options api.OnResolveOptions, callback func(api.OnResolveArgs) (api.OnResolveResult, error)) {
	options.Filter = "/\\.html$/"
	callback = func(args api.OnResolveArgs) (api.OnResolveResult, error) {
		path, _ := utils.FS.ConvertPath("./index.html")
		return api.OnResolveResult{
			Path: path,
		}, nil
	}
}
func loadHtml(options api.OnLoadOptions, callback func(api.OnLoadArgs) (api.OnLoadResult, error)) {
	path, _ := utils.FS.ConvertPath("./index.html")
	buf, _ := ioutil.ReadFile(path)
	html := string(buf)
	callback = func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		return api.OnLoadResult{
			Loader:   api.LoaderFile,
			Contents: &html,
		}, nil
	}
}
