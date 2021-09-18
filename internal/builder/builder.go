package builder

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"os"
)

// EsbuildStarter esbuild启动器 TODO:内存模式和文件模式切换
func EsbuildStarter() *api.BuildResult {
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{"D:/learn/espack/internal/builder/input.jsx"},
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
