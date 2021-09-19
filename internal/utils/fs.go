package utils

import (
	cPath "path"
	"path/filepath"
	"strings"
)

// 无实际意义，仅作为区分使用
type fs struct {
}

// ConvertPath 文件路径转换
func (f fs) ConvertPath(path string) (string, error) {
	if Env.Dev() {
		if strings.HasPrefix(path, "/") {
			path = cPath.Clean(path)[1:]
		}
		return filepath.Abs("./internal/case/" + path)
	}

	return filepath.Abs(path)
}
