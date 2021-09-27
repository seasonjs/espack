package mod

//npm 元数据获取 https://github.com/npm/registry/blob/master/docs/REGISTRY-API.md
//元数据字段解析  https://github.com/npm/registry/blob/master/docs/responses/package-metadata.md
//TODO:支持按照版本读取读取并下载 需要调研版本解析规范： https://semver.org/
// 全局的数据是否要放到sqlite中？

// 转换逻辑 package.js → js.mod → js.sum → 读取存在本地的sqlite 数据 ->没有数据则全量拉取
//                               ↓
//								 	[] mod_info //将所有层级平铺
//                               ↓
//									 mod_info //只会显示当前层级的require

//=====================espack get=================================

type mod struct {
}

func NewMod() *mod {
	return &mod{}
}

// UnTarMod UnTarMod 解压模块
func (m *mod) UnTarMod(input, output string) (err error) {

	return nil
}

func (m *mod) AnalyzeDependencies() {

}

// DownLoadDependencies 下载依赖 TODO:开十个协程限制数量，过多的协程可能使得效率反而下降
func (m *mod) DownLoadDependencies() {

}
