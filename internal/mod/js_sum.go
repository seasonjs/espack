package mod

import (
	"fmt"
	"seasonjs/espack/internal/utils"
	"strings"
)

// js.mod 的meta数据
type jsSum struct {
	modList    []*modInfo          //modList cache 为空即可达到边界
	knownList  map[string]*modInfo //拍平后的结构 name@version modInfo
	unKnowList map[string]string   //name@version string
}

func NewJsSum() *jsSum {
	return &jsSum{}
}

// FetchRootModMeta 获取根js.mod的meta数据
func (j *jsSum) FetchRootModMeta() *jsSum {
	//TODO:等测试结束移除该接口
	path, _ := utils.FS.ConvertPath("../case/js.mod")
	require := NewJsMod().ReadFile(path).Require
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
func (j *jsSum) PrunerAndFetcher() {
	if len(j.unKnowList) > 0 {
		for nv, version := range j.unKnowList {
			name := strings.Trim(nv, "@"+version)
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
}

func (j *jsSum) GetJsSum() jsSum {
	return *j
}
