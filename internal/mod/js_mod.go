package mod

import (
	"bufio"
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/pkg/errors"
	"os"
	"regexp"
	"seasonjs/espack/internal/utils"
	"strings"
)

// AnalyzeDependencies 分析Dependencies
// espack的实现方案：
// 1.读取package.json 然后解析dependencies字段+devDependencies字段+peerDependencies字段,去重 生成mod
// 2.  TODO: 递归第一步如果发生嵌套，并发请求npm元数据每个依赖生成js.mod解析去重,合并到mod上
// 3.  TODO: 是否要生成go mod like的js.mod 文件？如果要生成的话就需要在做js.mode解析
// 3.5 TODO: 如果要生成js.mod那么还需要node_model吗？是否可以通过改写esbuild 直接从全局缓存中拿文件？
// 4.根据分析出的依赖生成结构并发下载
// 依赖树(js.mod 或者package.json 以js.mod为优先，如果有js.mod就不会去扫描package.json)
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
// js.sum
// 即生成类package.json 或者yarn.lock目录

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

// 第一阶段只解析这些信息
type jsMod struct {
	Module  string
	Version string
	Main    string
	Target  api.Target
	Require map[string]string
}

func NewJsMod() *jsMod {
	return &jsMod{}
}

// ReadFile 读取JsMod
func (j *jsMod) ReadFile(p ...string) *jsMod {
	//判断是否存在js.mod文件
	path := ""
	if p != nil && len(p[0]) > 0 {
		path = p[0]
	} else {
		path, _ = utils.FS.ConvertPath("./js.mod")
	}
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			//TODO:如果没有js.mod就读取package.json
			utils.Err.LogAndExit(errors.Wrap(err, "js.mod 文件不存在！"))
		} else {
			utils.Err.LogAndExit(errors.Wrap(err, "js.mod 读取失败"))
		}
	}
	file, _ := os.Open(path)
	// 关闭文件
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	reg, _ := regexp.Compile(`//.*$`)
	spaceReg, _ := regexp.Compile(`\s+`)
	var requireTokenStack []string
	//TODO: 完善js.mod读取 需要重写io.Reader接口
	for i := 1; fileScanner.Scan(); i++ {
		// 以//开头视为注释 空行和注释不读取
		if strings.HasPrefix(fileScanner.Text(), "//") || fileScanner.Text() == "" {
			continue
		}
		line := fileScanner.Text()
		//正则去除剩下的//注释，trim掉空格
		s := reg.ReplaceAllString(line, "")
		line = strings.TrimSpace(s)
		la := spaceReg.Split(line, -1)
		//TODO 改写成映射的格式
		if la[0] == "module" {
			if j.Module != "" {
				fmt.Println("第", i, "行,出现重复字段")
			}
			if len(la) == 2 {
				j.Module = la[1]
			} else {
				fmt.Println("第", i, "行格式错误")
			}
			continue
		}
		//TODO:是否要检查是合法版本
		if la[0] == "version" {
			if j.Version != "" {
				fmt.Println("第", i, "行,出现重复字段")
			}
			if len(la) == 2 {
				j.Version = la[1]
			} else {
				fmt.Println("第", i, "行格式错误")
			}
			continue
		}
		if la[0] == "main" {
			if j.Main != "" {
				fmt.Println("第", i, "行,出现重复字段")
			}
			if len(la) == 2 {
				j.Main = la[1]
			} else {
				fmt.Println("第", i, "行格式错误")
			}
			continue
		}
		// TODO 是否需要进行扩展支持多编译版本？
		if la[0] == "target" {
			if j.Target != api.DefaultTarget {
				fmt.Println("第", i, "行,出现重复字段")
			}
			if len(la) == 2 {
				switch la[1] {
				case "ES5":
					j.Target = api.ES5
				case "ES2015":
					j.Target = api.ES2015
				case "ES2016":
					j.Target = api.ES2016
				case "ES2017":
					j.Target = api.ES2017
				case "ES2018":
					j.Target = api.ES2018
				case "ES2019":
					j.Target = api.ES2019
				case "ES2020":
					j.Target = api.ES2020
				case "ES2021":
					j.Target = api.ES2021
				case "ESNext":
					j.Target = api.ESNext
				default:
					j.Target = api.DefaultTarget
					fmt.Println("错误的target")
				}
			} else {
				fmt.Println("第", i, "行格式错误")
			}
			continue
		}

		if la[0] == "require" {
			if len(la) == 2 {
				continue
			}
			if len(la) == 1 {
				continue
			}
			fmt.Println("第", i, "行格式错误")
			continue
		}
		if la[0] == "(" {
			fmt.Println("第", i, "行格式错误")
			continue
		}
		fmt.Println("第", i, "行格式错误")

	}
	//var separatorTrack []string
	//for len(tokenStack) > 0 {
	//	x := tokenStack[0]
	//	tokenStack = tokenStack[1:]
	//	// 如果操作符栈为空，则说明是开始
	//	if len(separatorTrack) == 0 {
	//		separatorTrack = append(separatorTrack, x)
	//		continue
	//	}
	//	// 如果操作符栈为1，则说明是代码块内容
	//	if len(separatorTrack) == 1 {
	//		separatorTrack = append(separatorTrack, x)
	//	}
	//	// 如果操作符栈大于1，则说明是出现了问题
	//	if len(separatorTrack) > 1 {
	//		fmt.Println("解析失败 ：）")
	//	}
	//	fmt.Println(x)
	//}
	//fmt.Println(tokenStack)
	if err != nil {
		fmt.Println("读取文件失败:", err)
	}
	return j
}

//
type modInfo struct {
}

func NewModInfo() *modInfo {
	return &modInfo{}
}

// GetJsMod 获取模块的信息
func (m *modInfo) GetJsMod(name string) {

}
