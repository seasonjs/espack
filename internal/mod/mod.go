package mod

//npm 元数据获取 https://github.com/npm/registry/blob/master/docs/REGISTRY-API.md
//元数据字段解析  https://github.com/npm/registry/blob/master/docs/responses/package-metadata.md
//TODO:支持按照版本读取读取并下载 需要调研版本解析规范： https://semver.org/
//=====================espack get=================================
type mod struct {
}

func NewMod() *mod {
	return &mod{}
}

// UntarMod 解压模块
func (m *mod) UnzipMod(input, output string) (err error) {

	return nil
}

func (m *mod) AnalyzeDependencies() {

}

// DownLoadDependencies 下载依赖 TODO:开十个协程限制数量，过多的协程可能使得效率反而下降
// npm 包下载代码仓库 1. https://github.com/npm/pacote
// 最基本的包下载 https://github.com/npm/npm-registry-fetch
//https://registry.npmjs.org/${packageName}/-/${packageName}-${version}.tgz

func (m *mod) DownLoadDependencies() {

}
