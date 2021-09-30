package mod

import (
	"encoding/json"
	"fmt"
	"net/http"
	"seasonjs/espack/internal/logger"
	"seasonjs/espack/internal/utils"
	"time"
)

// TODO:需要同步配置npm的：）
// const pkg = require('./package.json')
// module.exports = {
//  log: require('./silentlog.js'),
//  maxSockets: 12,
//  method: 'GET',
//  registry: 'https://registry.npmjs.org/',
//  timeout: 5 * 60 * 1000, // 5 minutes
//  strictSSL: true,
//  noProxy: process.env.NOPROXY,
//  userAgent: `${pkg.name
//    }@${
//      pkg.version
//    }/node@${
//      process.version
//    }+${
//      process.arch
//    } (${
//      process.platform
//    })`,
// }
var _httpCli = &http.Client{
	Timeout: time.Duration(15) * time.Second,
	Transport: &http.Transport{
		MaxIdleConnsPerHost:   1,
		MaxConnsPerHost:       1,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

// 从npm远程获取的mod元数据，只解析需要解析的字段，减少解析耗时
type modMeta struct {
	Dist struct {
		Integrity    string `json:"integrity"`
		Shasum       string `json:"shasum"`
		Tarball      string `json:"tarball"`
		FileCount    int    `json:"fileCount"`
		UnpackedSize int    `json:"unpackedSize"`
		NpmSignature string `json:"npm-signature"`
	} `json:"dist"`
	Dependencies map[string]string `json:"dependencies"`
	//DevDependencies map[string]string `json:"devDependencies"`
}
type disMeta struct {
	Latest string `json:"latest"`
}
type modInfo struct {
	Name      string //包名
	Version   string //版本
	TarBall   string //下载地址
	Shasum    string
	integrity string
	Registry  string            // 元数据链接
	Require   map[string]string //当前包下的依赖
}

func NewModInfo(name, version string) *modInfo {
	// 需要从配置中获取如果没有则使用默认的
	tarball := fmt.Sprintf("https://registry.npmjs.org/%s/-/%s-%s.tgz", name, name, version)
	registry := fmt.Sprintf("https://registry.npmjs.org/%s/%s", name, version)
	return &modInfo{
		Name:     name,
		Version:  version,
		TarBall:  tarball,
		Registry: registry,
		Require:  make(map[string]string), //判断是否为空即可知道是否到达树顶
	}
}

// npm 包下载代码仓库 1. https://github.com/npm/pacote
// 最基本的包下载 https://github.com/npm/npm-registry-fetch
//https://registry.npmjs.org/${packageName}/-/${packageName}-${version}.tgz
// curl --location --request GET 'https://registry.npmjs.org/${packeageName}' \
// --header 'Accept: application/vnd.npm.install-v1+json'
//curl --location --request GET 'https://registry.npmjs.org/-/package/yarn/dist-tags'
//

//FetchModInfo TODO: 优化为异步拉取？考虑失败问题：）
func (i *modInfo) FetchModInfo() *modInfo {
	req, _ := http.NewRequest("GET", i.Registry, nil)
	// 通过设置请求头缩小元数据量
	req.Header.Set("Accept", "application/vnd.npm.install-v1+json; q=1.0, application/json; q=0.8, */*")
	rsp, err := (_httpCli).Do(req)
	if err != nil {
		logger.Warn("获取包的元数据信息失败，%s", err.Error())
		return i
	}
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusOK {
		decoder := json.NewDecoder(rsp.Body)
		//通过mm解析json
		mm := modMeta{}
		err = decoder.Decode(&mm)
		if err != nil {
			//TODO:重试，不要退出
			logger.Warn("解析包%s的元数据失败", i.Registry)
			return i
		}
		i.Shasum = mm.Dist.Shasum
		i.integrity = mm.Dist.Integrity
		if len(mm.Dependencies) > 0 {
			for name, version := range mm.Dependencies {
				// npm 的package.json 存在版本方言 ：）
				v := utils.Version.FindVersionStr(version)
				//最懒方案 ^^ 如果存在版本方言则选择拉取最新的版本
				if len(v[""]) == 0 {
					i.Require[name] = i.FetchLastVersion(name)
					continue
				}
				i.Require[name] = v[""]

			}
		}
	} else {
		//TODO:retry
	}
	//if len(mm.DevDependencies) > 0 {
	//	for name, version := range mm.DevDependencies {
	//		// npm 的package.json 存在版本方言 ：）
	//		v := utils.Version.FindVersionStr(version)
	//		//最懒方案 ^^ 如果存在版本方言则选择拉取最新的版本
	//		if len(v[""]) == 0 {
	//			i.Require[name] = i.FetchLastVersion(name)
	//			continue
	//		}
	//		i.Require[name] = v[""]
	//	}
	//}
	return i
}

func (i *modInfo) FetchLastVersion(name string) string {
	tagsUrl := fmt.Sprintf("https://registry.npmjs.org/-/package/%s/dist-tags", name)
	req, _ := http.NewRequest("GET", tagsUrl, nil)
	// 通过设置请求头缩小元数据量
	req.Header.Set("Accept", "application/vnd.npm.install-v1+json; q=1.0, application/json; q=0.8, */*")
	//TODO 是否不需要每次都初始化Client
	rsp, err := (&http.Client{}).Do(req)
	if err != nil {
		logger.Warn("获取包最新版本失败，%s", err.Error())
	}
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusOK {
		decoder := json.NewDecoder(rsp.Body)
		rsp.Header.Get("")
		//通过mm解析json
		dm := disMeta{}
		err = decoder.Decode(&dm)
		if err != nil {
			//TODO:重试，不要退出
			logger.Warn("包最新版本%s的信息解析失败", i.Registry)
		}
		return dm.Latest
	} else {
		//TODO： retry
	}
	return ""
}

//TODO: 是否要进行 版本方言判断？ https://docs.npmjs.com/cli/v7/configuring-npm/package-json#dependencies
func isDialectVersion(ver string) bool {
	//只解决 latest, * 不会过多解决
	if ver == "latest" {
		return true
	}
	if ver == "*" {
		return true
	}
	// 巴拉巴拉... :)
	return false
}
