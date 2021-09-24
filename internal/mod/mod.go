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
// espack的实现方案：
// 1.读取package.json 然后解析dependencies字段+devDependencies字段+peerDependencies字段,去重 生成mod
// 2.  TODO: 递归第一步如果发生嵌套，并发请求npm元数据每个依赖生成modInfo解析去重,合并到mod上
// 3.  TODO: 是否要生成go mod like的js.mod 文件？如果要生成的话就需要在做js.mode解析
// 3.5 TODO: 如果要生成js.mod那么还需要node_model吗？是否可以通过改写esbuild 直接从全局缓存中拿文件？
// 4.根据分析出的依赖生成结构并发下载
// 依赖树
// root -> a(v1) , b(v2)
// a(v1)->b(v1)
//       |_____ b(v1)->c(v1)
// b(v2)->c(v2)
// c(v2)->b(v2)
// 包结构（全局共享一个依赖目录）直接拍平获取即只在
// root
// +-- a_v1
// +-- b_v1
// +-- c_v1
// +-- b_v2
// +-- c_v2
//

// npm yarn 的实现方案与问题
// 讨论issue： https://github.com/npm/cli/issues/1597#issuecomment-667639545
// 源码：https://github.com/npm/cli/blob/latest/node_modules/%40npmcli/arborist/lib/arborist/reify.js#L799
// Example dep graph:
// root -> (a, c)
// a -> BUNDLE(b)
// b -> c
// c -> b
//
// package tree:
// root
// +-- a
// |   +-- b(1)
// |   +-- c(1)
// +-- b(2)
// +-- c(2)
// 1. mark everything that's shadowed by anything in the bundle.  This
//    marks b(2) and c(2).
// 2. anything with edgesIn from outside the set, mark not-extraneous,
//    remove from set.  This unmarks c(2).
// 3. continue until no change
// 4. remove everything in the set from the tree.  b(2) is pruned
// lib -> (a@1.x) BUNDLE(a@1.2.3 (b@1.2.3))
// a@1.2.3 -> (b@1.2.3)
// a@1.3.0 -> (b@2)
// b@1.2.3 -> ()
// b@2 -> (c@2)
//
// root
// +-- lib
// |   +-- a@1.2.3
// |   +-- b@1.2.3
// +-- b@2 <-- shadowed, now extraneous
// +-- c@2 <-- also shadowed, because only dependent is shadowed

func (m *mod) AnalyzeDependencies() {

}

// DownLoadDependencies 下载依赖 TODO:开十个协程限制数量，过多的协程可能使得效率反而下降
// npm 包下载代码仓库 1. https://github.com/npm/pacote
// 最基本的包下载 https://github.com/npm/npm-registry-fetch

//https://registry.npmjs.org/${packageName}/-/${packageName}-${version}.tgz
func (m *mod) DownLoadDependencies() {

}
