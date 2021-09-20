package builder

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"os"
	"seasonjs/espack/internal/utils"
)

// EsbuildStarter esbuild启动器 TODO:内存模式和文件模式切换 是否要对esbuilder的 ast进行提取？
func EsbuildStarter() *api.BuildResult {
	//ap, _ := utils.FS.ConvertPath("./input.jsx")
	//htm, _ := utils.FS.ConvertPath("./index.html")
	prefix := ""
	if utils.Env.Dev() {
		prefix = "./internal/case"
	}

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{prefix + "/index.jsx"},
		//Outfile:     "output.js",
		Loader: map[string]api.Loader{
			".html": api.LoaderFile,
			".svg":  api.LoaderDataURL,
		},
		Outdir:   prefix + "/dist",
		Bundle:   true,
		Write:    false,
		LogLevel: api.LogLevelInfo,
		Target:   api.ES2016,
		//Plugins:  []api.Plugin{htmlBuilder.NewHtmlPlugin()},
	})
	fmt.Printf("%d errors and %d warnings\n",
		len(result.Errors), len(result.Warnings))
	if len(result.Errors) > 0 {
		os.Exit(1)
	}
	return &result
}
