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

// UnzipMod 解压模块
func (m *mod) UnzipMod(input, output string) (err error) {

	return nil
}

type modInfo struct {
}

func NewModInfo() *modInfo {
	return &modInfo{}
}

// GetModInfo 获取单个模块的信息
func (m *modInfo) GetModInfo(name string) {

}

// AnalyzeDependencies 分析Dependencies
// 1.读取package.json 然后解析dependencies字段+devDependencies字段+peerDependencies字段,去重 生成mod
// 2.  TODO: 递归第一步如果发生嵌套，并发请求npm元数据每个依赖生成modInfo解析去重,合并到mod上
// 3.  TODO: 是否要生成go mod like的js.mod 文件？如果要生成的话就需要在做js.mode解析
// 3.5 TODO: 如果要生成js.mod那么还需要node_model吗？是否可以通过改写esbuild 直接从全局缓存中拿文件？
// 4.根据分析出的依赖生成结构并发下载
func (m *mod) AnalyzeDependencies() {

}

// DownLoadDependencies 下载依赖TODO:开十个协程限制数量，过多的协程可能使得效率反而下降
func (m *mod) DownLoadDependencies() {

}
