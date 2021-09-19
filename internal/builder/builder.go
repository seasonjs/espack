package builder

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"os"
	"seasonjs/espack/internal/utils"
)

// EsbuildStarter esbuild启动器 TODO:内存模式和文件模式切换
func EsbuildStarter() *api.BuildResult {
	ap, _ := utils.FS.ConvertPath("./input.jsx")
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{ap},
		Outfile:     "output.js",
		Bundle:      true,
		Write:       false,
		LogLevel:    api.LogLevelInfo,
	})
	fmt.Printf("%d errors and %d warnings\n",
		len(result.Errors), len(result.Warnings))
	if len(result.Errors) > 0 {
		os.Exit(1)
	}
	return &result
}
