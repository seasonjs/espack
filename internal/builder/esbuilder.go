package builder

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"os"
	"seasonjs/espack/internal/config"
)

type esbuild struct {
	options       api.BuildOptions
	configuration config.Configuration
}

func NewEsBuilder(configuration config.Configuration) *esbuild {
	return &esbuild{
		api.BuildOptions{},
		configuration,
	}
}

// EsbuildStarter esbuild启动器 TODO:内存模式和文件模式切换 是否要对esbuild的 ast进行提取？
func (e *esbuild) EsbuildStarter() *api.BuildResult {
	//prefix := ""
	//if utils.Env.Dev() {
	//	prefix = "./internal/case"
	//}
	////TODO: 对esbuild的log进行处理
	//result := api.Build(api.BuildOptions{
	//	EntryPoints: []string{prefix + "/index.jsx"},
	//	//Outfile:     "output.js",
	//	Loader: map[string]api.Loader{
	//		".html": api.LoaderFile,
	//		".svg":  api.LoaderDataURL,
	//	},
	//	Outdir:   prefix + "/dist",
	//	Bundle:   true,
	//	Write:    false,
	//	LogLevel: api.LogLevelInfo,
	//	Target:   api.ES2016,
	//	//Plugins:  []api.Plugin{htmlPlugin.NewHtmlPlugin()},
	//})
	result := api.Build(e.options)
	fmt.Printf("%d errors and %d warnings\n",
		len(result.Errors), len(result.Warnings))
	if len(result.Errors) > 0 {
		fmt.Println(result)
		os.Exit(1)
	}
	return &result
}
