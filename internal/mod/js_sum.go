package mod

import (
	"fmt"
	"github.com/seasonjs/espack/internal/utils"
	"io"
	"os"
	"strings"
)

type JsSumAbs interface {
	FetchRootModMeta(map[string]string) *jsSum
	PrunerAndFetcher(path ...string) *jsSum
	GetJsSum() *jsSum
	WriteFile() *jsSum
}

// js.mod 的meta数据
type jsSum struct {
	modList    []*modInfo          //modList cache 为空即可达到边界
	knownList  map[string]*modInfo //拍平后的结构 name@version modInfo
	unKnowList map[string]string   //name@version string
}

func NewJsSum() *jsSum {
	return &jsSum{
		knownList:  make(map[string]*modInfo),
		unKnowList: make(map[string]string),
	}
}

// FetchRootModMeta 获取根js.mod的meta数据
func (j *jsSum) FetchRootModMeta(require map[string]string) *jsSum {
	for name, version := range require {
		key := fmt.Sprintf("%s@%s", name, version)
		////如果有这个包则不需要在获取它的信息
		//if _, ok := j.knownList[key]; ok {
		//
		//}

		mi := NewModInfo(name, version).FetchModInfo()
		j.knownList[key] = mi
		j.modList = append(j.modList, mi)
	}
	//需要等待全部被获得才能得到准确的列表
	for _, info := range j.modList {
		for name, version := range info.Require {
			subKey := fmt.Sprintf("%s@%s", name, version)
			//如果不存在这个key,则需要继续请求元数据
			if _, ok := j.knownList[subKey]; !ok {
				//利用map特性去重
				j.unKnowList[subKey] = version
			}
		}
	}
	//释放空间
	j.modList = nil
	return j
}

// PrunerAndFetcher 剪枝&与递归请求
func (j *jsSum) PrunerAndFetcher() *jsSum {
	if len(j.unKnowList) > 0 {
		for nv, version := range j.unKnowList {
			//@type/string@1.2.0 trim type/string
			//strings.TrimRight @type/string
			name := strings.TrimRight(nv, "@"+version)
			mi := NewModInfo(name, version).FetchModInfo()
			j.knownList[nv] = mi
			j.modList = append(j.modList, mi)
		}
		//清空
		j.unKnowList = make(map[string]string)
		//需要等待全部被获取，才能对列表进行剪枝
		for _, info := range j.modList {
			for name, version := range info.Require {
				subKey := fmt.Sprintf("%s@%s", name, version)
				//如果不存在这个key,则需要继续请求元数据
				if _, ok := j.knownList[subKey]; !ok {
					//利用map特性去重
					j.unKnowList[subKey] = version
				}
			}
		}
		//释放空间
		j.modList = nil
		j.PrunerAndFetcher()
	}
	return j
}

func (j *jsSum) WriteFile(path ...string) *jsSum {
	var filePath string
	if len(path) == 1 && len(path[0]) > 0 {
		filePath = path[0]
	} else {
		filePath, _ = utils.FS.ConvertPath("./js.sum")
	}
	// 覆盖写入，如果不存在就创建
	f, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE, 0666) //打开文件
	defer f.Close()
	for s, info := range j.knownList {
		_, err = io.WriteString(f, fmt.Sprintf("%s %s %s %s\n", s, info.TarBall, info.Shasum, info.integrity)) //写入文件(字符串)
		for name, version := range info.Require {
			_, err = io.WriteString(f, fmt.Sprintf("\t%s %s\n", name, version))
		}
		_, err = io.WriteString(f, fmt.Sprintf("\n"))
	}
	if err != nil {
		panic(err)
	}
	return j
}

// WriteDB 将元数据写入数据库 TODO
func (j *jsSum) WriteDB(path ...string) *jsSum {

	return j
}

func (j *jsSum) GetJsSum() jsSum {
	return *j
}
