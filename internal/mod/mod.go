package mod

import (
	"github.com/seasonjs/espack/internal/logger"
	"github.com/seasonjs/espack/internal/utils"
	"net/http"
	"path/filepath"
)

//npm 元数据获取 https://github.com/npm/registry/blob/master/docs/REGISTRY-API.md
//元数据字段解析  https://github.com/npm/registry/blob/master/docs/responses/package-metadata.md
//TODO:支持按照版本读取读取并下载 需要调研版本解析规范： https://semver.org/
// 全局的数据是否要放到sqlite中？

// 转换逻辑 package.js → js.mod → js.sum → 读取存在本地的sqlite 数据  → 没有数据则全量拉取 →下载所有tarball →全部解压
//                               ↓
//								 	[] mod_info //将所有层级平铺
//                               ↓
//									 mod_info //只会显示当前层级的require

//=====================espack get=================================

type mod struct {
	TarBall map[string]string
}

func NewMod() *mod {
	return &mod{}
}

// AnalyzeDependencies TODO 暂且跳过包分析全量下载生成
func (m *mod) AnalyzeDependencies(p ...string) *mod {
	path := ""
	jm := NewJsMod()
	var mf *jsMod
	if p != nil && len(p) > 0 {
		path = p[0]
		mf = jm.ReadFile(path)
	} else {
		mf = jm.ReadFile()
	}
	//读取mod文件
	require := mf.Require
	//递归mod文件
	jsPtr := NewJsSum().FetchRootModMeta(require).PrunerAndFetcher()
	var js jsSum
	if p != nil && len(p) > 1 {
		jsPtr.WriteFile(p[1])
	} else {
		jsPtr.WriteFile()
	}
	js = jsPtr.GetJsSum()
	m.TarBall = make(map[string]string, len(js.knownList))
	for _, info := range js.knownList {
		////现在是为了esbuild的兼容性 TODO:将版本注入
		//m.TarBall = append(m.TarBall, info.TarBall)
		m.TarBall[info.Name] = info.TarBall
	}
	return m
}

// DownLoadDependencies 下载依赖 TODO:开十个协程限制数量，过多的协程可能使得效率反而下降
func (m *mod) DownLoadDependencies(p ...string) {

	path := ""
	if p != nil && len(p) > 0 {
		path = p[0]
	} else {
		//TODO 目前为了测试esbuild方案可行，所以不移动代码到根目录
		path, _ = utils.FS.ConvertPath("./node_modules/")
	}
	for name, url := range m.TarBall {
		logger.Info("开始下载%s,资源地址%s", name, url)
		resp, _ := http.Get(url)
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			err := DefaultTarGz().IOUnTarGz(resp.Body, filepath.Join(path, name))
			if err != nil {
				logger.Warn("解压%s包文件失败", name)
			}
			logger.Info("%s下载并解压成功", name)
		} else {
			//...
		}
	}

}
