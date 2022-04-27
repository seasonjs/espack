package utils

import (
	"os"
	cPath "path"
	"path/filepath"
	"runtime"
	"strings"
)

// 无实际意义，仅作为区分使用
type fs struct {
}

// ConvertPath 文件路径转换
func (f fs) ConvertPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	if Env.Dev() {
		if strings.HasPrefix(path, "/") {
			path = cPath.Clean(path)[1:]
		}
		return filepath.Abs("./internal/case/" + path)
	}

	return filepath.Abs(path)
}
func (f fs) IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func (f fs) GetCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return cPath.Dir(filename)
}
