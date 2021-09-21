package utils

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"runtime"
)

// 无实际意义，仅作为区分使用
type env struct {
}

func newEnv() {
	//获取当前调用的文件的路径
	_, filename, _, _ := runtime.Caller(1)
	path := filepath.Join(filename, "../env/.env")
	e := godotenv.Load(path)
	if e != nil {
		Err.LogAndExit(e)
	}
}
func (e env) Dev() bool {
	return os.Getenv("ES_PACK_DEV_ENV") == "dev"
}

func (e env) Prod() bool {
	return os.Getenv("ES_PACK_DEV_ENV") == "prod"
}
