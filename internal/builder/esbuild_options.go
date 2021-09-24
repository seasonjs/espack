package builder

import (
	"github.com/evanw/esbuild/pkg/api"
)

// GetOptions 对配置信息进行转换
func (e *esbuild) GetOptions() *esbuild {
	var entryPointsAdvanced []api.EntryPoint
	for key, value := range e.configuration.Entry {
		entryPointsAdvanced = append(entryPointsAdvanced, api.EntryPoint{
			InputPath:  value,
			OutputPath: key,
		})
	}
	e.options = api.BuildOptions{
		EntryPointsAdvanced: entryPointsAdvanced,
		Outdir:              e.configuration.Output.Path,
		AbsWorkingDir:       e.configuration.Context,
		Loader: map[string]api.Loader{ //TODO:是否需要将Loader格式写死
			".html": api.LoaderFile,
			".svg":  api.LoaderDataURL,
		},
		//LogLevel: api.LogLevelDebug, //TODO: 自定义消息
		Bundle: true, //TODO:调研Bundle的问题是否要暴露出来让用户配置
	}

	return e
}
