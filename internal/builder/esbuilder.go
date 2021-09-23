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
	result := api.Build(e.options) //TODO: 对esbuild的log进行处理 错误输出到页面而不是终止执行
	fmt.Printf("%d errors and %d warnings\n",
		len(result.Errors), len(result.Warnings))
	if len(result.Errors) > 0 {
		fmt.Println(result)
		os.Exit(1)
	}
	return &result
}
